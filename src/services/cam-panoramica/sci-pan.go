package scipan

import (
	"errors"
	"fmt"

	"github.com/gustavolimam/control-access/src/components/buffer"
	"github.com/gustavolimam/control-access/src/components/camera"
	"github.com/gustavolimam/control-access/src/components/config"
	"github.com/gustavolimam/control-access/src/components/log"
	"github.com/gustavolimam/control-access/src/components/messages"
)

const (
	logService log.Service = "CAM-PANORAMICA"
)

var (
	errLateImg  = errors.New("Imagem do buffer fora do intervalo de tempo permitido (Atrasada)")
	errEarlyImg = errors.New("Imagem do buffer fora do intervalo de tempo permitido (Adiantada)")
	erroImgNull = errors.New("Imagem não encontrada no buffer")
)

// SciPan representa a estrutura do canal sci-pan
type SciPan struct {
	sdpCh  chan messages.PanReceive
	cam    *camera.Camera
	buffer *buffer.FrameBuffer
}

// New inicia um novo serviço do SCI-PAN
func New() *SciPan {
	log.Log(logService, "Serviço criado")
	return &SciPan{
		sdpCh: messages.GetSdpToPan(),
		cam: camera.New(
			logService,
			config.Config.PanCam.Address,
			messages.GetChanPan(),
			config.Config.PanCam.FrameRate,
			config.Config.PanCam.ImgQuality,
		),
		buffer: buffer.NewBuffer(defaults.BufferSize)}
}

// Run função resposanvel pelas principais chamadas das cameras
func (s *SciPan) Run() {
	log.Log(logService, "Serviço iniciado")

	go s.cam.SendFrames()

	//Fica recebendo os frames da panoramica, processando-as e salvando no buffer
	go func() {
		for {
			c := <-s.cam.ChFrame
			img := s.cam.ProcessaFrame(c)
			s.buffer.Add(img)
		}
	}()

	//Fica escutando o canal do SDP até receber uma imagem com placa
	//Quando chega alguma, busca no buffer da panoramica a imagem correspondente
	go func() {
		for {
			f := <-s.sdpCh
			var msg *messages.Msg
			if msg = s.buscaFrame(f); msg.Err == nil {
				msg.Err = s.ValidaTimeFrame(*msg, f)
			}

			if msg.Err != nil {
				log.Log(logService, fmt.Sprintf("Mensagem enviada ao SLP (ID: %4.d) com erro: %s", msg.ID, msg.Err.Error()))
			} else {
				log.Log(logService, fmt.Sprintf("Mensagem enviada ao SLP. ID: %4.d", msg.ID))
			}
			messages.SendPanToScd(*msg)
		}
	}()
}

//Procura no buffer a imagem com o timestamp mais proximo ao recebido
func (s *SciPan) buscaFrame(panRcv messages.PanReceive) *messages.Msg {
	var img *image.ImageStruct
	var err error

	if img, err = s.buffer.Frame(panRcv.Time); err != nil {
		return &messages.Msg{
			Err: err,
		}
	}
	return &messages.Msg{
		ID:    panRcv.ID,
		Frame: *img,
	}
}

//ValidaTimeFrame se o frame retornado pelo buffer está dentro do tempo maximo e mínimo
func (s *SciPan) ValidaTimeFrame(msg messages.Msg, f messages.PanReceive) error {
	bufferTime := f.Time
	timeZoom := msg.Frame.Time

	timeMax := timeZoom.Add(defaults.TimeMax)
	timeMin := timeZoom.Add(defaults.TimeMin)

	if bufferTime.After(timeMax) {
		return errEarlyImg
	}

	if bufferTime.Before(timeMin) {
		return errLateImg
	}

	return nil
}
