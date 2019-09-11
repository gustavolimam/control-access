package buffer

import (
	"sync"
	"time"

	"github.com/gustavolimam/control-access/src/components/config"
	"github.com/gustavolimam/control-access/src/components/messages"
)

// InfraBuffer define a estrutura do buffer circular de vídeo
type InfraBuffer struct {
	bufferMutex sync.Mutex
	bufferPlate messages.BufferPlate
}

// NewPlateBuffer cria um novo buffer
func NewPlateBuffer() *InfraBuffer {
	b := new(InfraBuffer)
	b.bufferPlate = make(messages.BufferPlate)
	b.bufferMutex = sync.Mutex{}
	return b
}

// FindPlateBuffer busca o índice do elemento que tem o timestamp mais próximo ao tempo t
func (b *InfraBuffer) FindPlateBuffer(plateBuf *messages.BufferPackage) bool {
	plate := plateBuf.PlateInfo.Placa

	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()
	// Verifica se a placa já se encontra no buffer, caso não esteja ela será salva
	if _, ok := b.bufferPlate[plate]; !ok {
		b.bufferPlate[plate] = *plateBuf
		return false
	} else {
		// Após identificar a placa no buffer, verificamos o tempo e caso esteja acima do tempo configurado,
		// a função retorna true, permitindo que a placa seja passada adiante para a então geração da multa.
		timestamp := time.Since(b.bufferPlate[plate].Frame.Time).Seconds()
		if timestamp > config.Config.Plate.PlateBufferSeconds {
			return true
		}
	}
	return false
}

// DeletaPlateBuffer -  Percorre todo o map e verifica se algum dos itens encontrados estão a mais de 10 minutos,
// caso ultrapasse os 10 minutos, o mesmo será deletado do map.
func (b *InfraBuffer) DeletaPlateBuffer() {
	for {
		b.bufferMutex.Lock()
		for k := range b.bufferPlate {
			timestamp := time.Since(b.bufferPlate[k].Frame.Time).Minutes()
			if timestamp > bufferTimeout {
				delete(b.bufferPlate, k)
			}
		}
		b.bufferMutex.Unlock()
		time.Sleep(bufferTimeoutLoop)
	}
}
