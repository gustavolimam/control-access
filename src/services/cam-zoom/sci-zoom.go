package scizoom

import (
	"fmt"

	"github.com/gustavolimam/control-access/src/components/camera"
	"github.com/gustavolimam/control-access/src/components/config"
	"github.com/gustavolimam/control-access/src/components/log"
	"github.com/gustavolimam/control-access/src/components/messages"
)

const (
	logService log.Service = "CAM-ZOOM"
)

var (
	frameID = 0
)

// SciZoom representa a estrutua do serviço SCI-ZOOM
type SciZoom struct {
	cam *camera.Camera
}

// New retorna uma estrutura do servço sci-zoom
func New() *SciZoom {
	log.Log(logService, "Serviço criado")
	return &SciZoom{cam: camera.New(
		"CAM-ZOOM",
		config.Config.ZoomCam.Address,
		messages.GetChanZoom(),
		config.Config.ZoomCam.FrameRate,
		config.Config.ZoomCam.ImgQuality)}
}

// Run realiza a função do serviço sci-zoom:
// 1. Recebe frames da camera zoom
// 2. Processa o frame atribuindo informações (modo noturno e timestamp)
// 3. Envia para o sdp
func (s *SciZoom) Run() {
	log.Log(logService, "Serviço iniciado")

	go s.cam.SendFrames()

	for {
		// Recebimento de frames da camera
		if frameID > defaults.IDMax {
			frameID = 1
		}

		c := <-s.cam.ChFrame
		img := s.cam.ProcessaFrame(c)

		frameID++

		messages.SendResultSdp(&image.ImageZoomID{ZoomID: frameID, Img: img})
		log.Log(logService, fmt.Sprintln("Imagem Zoom enviada ao SDP - IDZoom = ", frameID))
	}
}
