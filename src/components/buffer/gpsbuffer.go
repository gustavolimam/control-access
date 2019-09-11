package buffer

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gustavolimam/control-access/src/components/log"
	"github.com/gustavolimam/control-access/src/components/messages"
)

const (
	logService = "GPS-BUFFER"
)

// GpsBuffer define a estrutura do buffer circular de dados do gps
type GpsBuffer struct {
	bufferMutex sync.Mutex
	bufferGPS   messages.BufferGPS
}

// NewGPSBuffer cria um novo buffer
func NewGPSBuffer() *GpsBuffer {
	b := new(GpsBuffer)
	b.bufferGPS = make(messages.BufferGPS)
	b.bufferMutex = sync.Mutex{}
	return b
}

// AddGPSBuffer adiciona um novo pacote de GPS no buffer
func (b *GpsBuffer) AddGPSBuffer(gpsBuf *messages.SglPackage) {
	timestampGPSs := ajustaTempo(gpsBuf.Dados.Timestamp)

	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()
	// Verifica se os dados de gps já se encontra no buffer, caso não esteja eles serão salvos
	if _, ok := b.bufferGPS[timestampGPSs]; !ok {
		b.bufferGPS[timestampGPSs] = *gpsBuf
	} else {
		log.Log(logService, "Pacote já inserido no buffer")
	}
}

// FindGPSBuffer busca o índice do elemento que tem o timestamp mais próximo ao tempo t
func (b *GpsBuffer) FindGPSBuffer(msg messages.MsgSgl) messages.SglPackage {
	timestampGPS := ajustaTempo(msg.Tempo)
	err := errors.New("Pacote não encontrado")

	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()
	// Verifica se o timestamp passado é válido, ou seja, não é igual a zero
	if !timestampGPS.IsZero() {
		// Percorre toda a extensão do buffer em busca de algum dado do GPS que corresponda
		// ao horario passado como parâmetro, após isso concatena os dados extras e transmite
		// o resultado com a estrutura final do GPS
		for tempo, k := range b.bufferGPS {
			if tempo == timestampGPS {
				pacote := messages.SglPackage{
					ID:    msg.ID,
					Dados: k.Dados,
					Erro:  nil,
				}
				return pacote
			}
		}
	}
	// Caso não encontre nada irá retorna somente o erro
	return messages.SglPackage{Erro: err}
}

// DeletaGPSBuffer -  Percorre todo o map e verifica se algum dos itens encontrados estão a mais de 10 minutos,
// caso ultrapasse os 10 minutos, o mesmo será deletado do map.
func (b *GpsBuffer) DeletaGPSBuffer() {
	for {
		b.bufferMutex.Lock()
		for k := range b.bufferGPS {
			timestamp := time.Since(b.bufferGPS[k].Dados.GGAgps.Time).Minutes()
			if timestamp > bufferTimeout {
				delete(b.bufferGPS, k)
			}
		}
		b.bufferMutex.Unlock()
		time.Sleep(bufferTimeoutLoop)
	}
}

// ajustaTempo retorna o timestamp em string para definir o padrão usado no buffer do GPS
func ajustaTempo(tempo time.Time) time.Time {

	hora := strconv.Itoa(tempo.Hour())
	minuto := strconv.Itoa(tempo.Minute())
	segundo := strconv.Itoa(tempo.Second())

	tempoStr := []string{hora, minuto, segundo}

	tempoString := strings.Join(tempoStr, ":")

	tempoFinal, err := time.Parse("15:04:05", tempoString)
	if err != nil {
		log.Log(logService, "Erro no parse do timestamp - ", err)
		return time.Time{}
	}

	return tempoFinal
}
