package camera

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gustavolimam/control-access/src/components/log"
)

const (
	timeStartAttempts     = 5               // maximo de requisições para tentar obter tempo inicial da câmera
	sleepRetryConection   = time.Second     // Tempo de espera para a tentativa de nova conexao a API camera
	timeoutImg            = time.Second     // Timeout recebimento de imagem
	blockConnection       = 2 * time.Second // Timeout para conexão bloqueada sem recebimento de nova imagem
	timeoutConnectionSync = 1 * time.Second // Timeout do client de conexão para sincronização do relógio
)

// Camera representa a estrutura de uma câmera
type Camera struct {
	logService         log.Service
	Address            string  // Ip da câmera
	ChFrame            ChFrame // Canal que os frames serão enviados
	FrameRate          int
	StartTimestamp     time.Time // Tempo inicial da câmera
	LastFrameTimestamp time.Time // Tempo do último frame
	ImgQuality         int
}

// header contains the fields that preceed a given image.
type header struct {
	boundary      string
	motionEvent   int
	contentType   string
	contentLength int
}

// ChFrame representa o canal de frames entre a câmera e outro serviço
type ChFrame chan []byte

// New retorna uma estrutura de câmera
func New(logService log.Service, address string, chFrame ChFrame, frameRate int, imgQuality int) *Camera {
	log.Log(logService, "Nova camera instanciada: ", address)

	return &Camera{logService, address, chFrame,
		frameRate, time.Now(), time.Now(), imgQuality}
}

// SendFrames envia frames capturados no canal ChFrame
func (c *Camera) SendFrames() {

	c.syncTimeLoopInfinito()

	URL := fmt.Sprintf("http://%s/api/mjpegvideo.cgi?Quality=%d&FrameRate=%d",
		c.Address, c.ImgQuality, c.FrameRate)

	client := http.Client{}
	var response *http.Response
	var err error
	for {
		response, err = client.Get(URL)
		if err != nil {
			log.Log(c.logService, "Falha na conexão HTTP de vídeo com a câmera: ", err.Error())
			time.Sleep(sleepRetryConection)
		} else {
			log.Log(c.logService, "Conexao feita (", c.Address, ") com sucesso.")
			break
		}
	}

	var getJpegOK int32
	// Monitora a conexão de vídeo com a câmera. Caso a função getJpeg()
	// fique bloqueado por mais de 2s fecha o response.
	// Body para forçar o retorno de erro.
	go func() {
		for {
			atomic.StoreInt32(&getJpegOK, 0)
			time.Sleep(blockConnection)
			if atomic.LoadInt32(&getJpegOK) == 0 {
				if response != nil {
					response.Body.Close()
				}
			}
		}
	}()

	// nextImg, no qual o sci-nmet.go vai ficar olhando e tratando os dados
	// processa dados do HTTP
	go func() {
		for {
			if response != nil {
				img, err := getJpeg(response.Body)
				atomic.StoreInt32(&getJpegOK, 1)
				// Caso não esteja mais recebendo imagens, dorme por um segundo e tenta reconectar posteriormente
				if err == io.EOF {
					time.Sleep(timeoutImg)
				}
				if err != nil {
					log.Log(c.logService, "Falha na decodificação do vídeo mjpeg: ", err.Error())
					response.Body.Close()
					response, err = client.Get(URL)
					if err != nil {
						log.Log(c.logService, "Falha na conexão http de vídeo com a câmera: ", err.Error())
					}
				}
				if img != nil {
					c.ChFrame <- img
				}
			} else {
				var err error
				c.syncTimeLoopInfinito()
				if response, err = client.Get(URL); err != nil {
					atomic.StoreInt32(&getJpegOK, 1)
					// Caso a URL esteja inacessível, aguarda um segundo antes de tentar novamente
					if err != nil {
						log.Log(c.logService, "Falha na tentativa de reconexão HTTP de vídeo com a camera. Erro: ",
							err.Error())
						time.Sleep(sleepRetryConection)
					} else {
						log.Log(c.logService, "Conexao feita (", c.Address, ") com sucesso.")
						continue
					}
					time.Sleep(sleepRetryConection)
				}
			}
		}
	}()
}

// syncTime atualiza c.inicioCamera para cálculo correto de timestamp das imagens
func (c *Camera) SyncTime() error {
	var err error
	// Tenta sincronizar o horário, caso não consiga retorna erro.
	// Se obtiver sucesso, define o tempo de sincronismo
	var init time.Time
	init, err = c.GetStartTime()
	if err == nil {
		c.StartTimestamp = init
	} else {
		log.Log(c.logService, "Não foi possível sincronizar o horário da câmera: ", err.Error())
	}
	return err
}

// syncTimeLoopInfinito executa sincronização do horário indefinidamente até que
// a sincronização obtenha sucesso, se isso ocorrer retorna
func (c *Camera) syncTimeLoopInfinito() {
	for {
		if err := c.SyncTime(); err == nil {
			return
		}
		time.Sleep(sleepRetryConection)
	}
}

// GetStartTime retorna o horário que a câmera foi iniciada
// Esse horário é estimado requisitando o TempoLigado da câmera
// t = horaRequisição + TempoResposta/2 - TempoLigado
func (c *Camera) GetStartTime() (t time.Time, err error) {
	url := fmt.Sprintf("http://%s/api/config.cgi?TempoLigado", c.Address)
	client := http.Client{Timeout: timeoutConnectionSync}
	d := 10 * time.Second
	for i := 0; i < timeStartAttempts; i++ {
		getTime := time.Now()
		responseTime, err := client.Get(url)
		dt := time.Since(getTime)
		if err == nil {
			tl, err := ioutil.ReadAll(responseTime.Body)
			responseTime.Body.Close()
			if err == nil {
				if d > dt {
					d = dt
					par := strings.Split(string(tl), "=")
					if tempo, err := strconv.Atoi(strings.TrimSpace(par[1])); err == nil {
						t = getTime.Add(dt/2 - (time.Duration(tempo) * time.Millisecond))
					} else {
						return t, err
					}
				}
			} else {
				return t, err
			}
		} else {
			return t, err
		}
	}
	return t, nil
}

// extractComment retorna os comentários do arquivo jpg
func (c *Camera) ExtractComment(img []byte) (data map[string]string) {
	data = map[string]string{}
	ini := bytes.Index(img, []byte{0xFF, 0xFE})
	if ini >= 0 {
		comString := string(img[ini+4:])
		for _, campo := range strings.Split(comString, ";") {
			par := strings.Split(campo, "=")
			if len(par) == 2 {
				data[par[0]] = par[1]
			}
		}
	}
	return
}

// processaTimestamp calcula o timestamp a partir do TempoCaptura, valida e sincroniza a base de tempo se necessário
func (c *Camera) ProcessaTimestamp(tempoCaptura uint64) time.Time {
	timeZero := time.Time{}
	frameTimestamp := c.StartTimestamp.Add(time.Duration(tempoCaptura) * time.Millisecond)
	// verifica se o timestamp é válido
	if frameTimestamp.Before(c.LastFrameTimestamp) {
		log.Log(c.logService, "Erro na verificação do timestamp - erro : Frame Anterior ", c.LastFrameTimestamp, "Valor do FrameTimestamp: ", frameTimestamp)
		c.syncTimeLoopInfinito()
		frameTimestamp = c.StartTimestamp.Add(time.Duration(tempoCaptura) * time.Millisecond)
	}
	/*if time.Since(frameTimestamp) > time.Second {
		log.Log(c.logService, "Erro na verificação do timestamp - erro : Frame Anterior ", c.LastFrameTimestamp, "Valor do FrameTimestamp: ", frameTimestamp)
		c.syncTimeLoopInfinito()
		frameTimestamp = c.StartTimestamp.Add(time.Duration(tempoCaptura) * time.Millisecond)
	}*/
	if c.StartTimestamp == timeZero {
		log.Log(c.logService, "Erro na verificação do timestamp - erro : Frame Anterior ", c.LastFrameTimestamp, "Valor do FrameTimestamp: ", frameTimestamp)
		c.syncTimeLoopInfinito()
		frameTimestamp = c.StartTimestamp.Add(time.Duration(tempoCaptura) * time.Millisecond)
	}

	c.LastFrameTimestamp = frameTimestamp
	return frameTimestamp
}

func (c *Camera) ProcessaFrame(img []byte) *image.ImageStruct {
	infoMap := c.ExtractComment(img)
	tempoCaptura, _ := strconv.ParseUint(strings.TrimSpace(infoMap["TempoCaptura"]), 10, 64)
	timestamp := c.ProcessaTimestamp(tempoCaptura)

	isNight, _ := strconv.ParseUint(strings.TrimSpace(infoMap["SituacaoDayNight"]), 10, 64)

	night := 0

	if isNight == 2 {
		night = 1
	}

	return &image.ImageStruct{Image: img, Time: timestamp, IsNightMode: night}
}

// getJpeg returns the next Image found in the MJPEG stream.
func getJpeg(inReader io.Reader) (img []byte, err error) {

	// read header
	h, err := readHeader(inReader)
	if err != nil {
		return nil, err
	}
	data := make([]byte, h.contentLength)
	switch h.contentType {
	case "image/jpeg":
		var s int
		for {
			n, err := inReader.Read(data[s:])
			s += n
			if err != nil {
				return nil, err
			}
			if s >= h.contentLength {
				return data, nil
			}
		}
	}
	return nil, errors.New("Unknown error")
}

// readHeader is used to find and return a correct MJPEG header.
func readHeader(inReader io.Reader) (h *header, outErr error) {
	// search for boundary
	data := make([]byte, 2)
	for {
		n, err := inReader.Read(data)
		switch {
		case err != nil:
			return nil, err
		case n == 1 && ((data[0] == 0x0A) || (data[0] == 0x0D)):
			continue
		case n < 2:
			return nil, fmt.Errorf("Not enough data available (2 needed - got %v: %X).", n, data[0])
		}
		if data[0] == '-' && data[1] == '-' {
			break
		}
	}

	// populate header
	h = new(header)
	h.boundary, outErr = readString(inReader, '\n')
	if outErr != nil {
		return nil, outErr
	}

	for {
		line, err := readString(inReader, '\n')
		if err != nil {
			return nil, err
		}

		if line == "" {
			break
		}

		kv := strings.Split(line, ": ")
		if len(kv) != 2 {
			return nil, errors.New("Not a valid key/value pair.")
		}

		switch kv[0] {
		case "Motion-Event":
			h.motionEvent, outErr = strconv.Atoi(kv[1])
			if outErr != nil {
				return nil, outErr
			}
		case "Content-Type":
			h.contentType = kv[1]
		case "Content-Length":
			h.contentLength, outErr = strconv.Atoi(kv[1])
			if outErr != nil {
				return nil, outErr
			}
		}
	}
	return h, nil
}

// readString tries to read a string until delim is found. The delim byte
// won't be returned.
func readString(inReader io.Reader, delim byte) (string, error) {
	var b = make([]byte, 1)
	buffer := bytes.NewBuffer(nil)
	for {
		n, err := inReader.Read(b)
		if err != nil || n < 1 {
			return "", err
		}
		if b[0] == delim {
			return strings.TrimSpace(buffer.String()), nil
		}
		buffer.Write(b)
	}
}
