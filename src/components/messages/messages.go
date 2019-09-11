package messages

import (
	"time"
)

// SglPackage representa a estrutura do pacote referente ao serviço de GPS
type SglPackage struct {
	ID    int
	Dados gps.GPSPackage
	Erro  error
}

type MsgSgl struct {
	ID    int
	Tempo time.Time
}

// Msg representa a estrutura do canal para comunicação entre scizoom e sdp
type Msg struct {
	ZoomID int
	ID     int
	Frame  image.ImageStruct
	Err    error
}

// PanReceive representa a estrutura do canal para comunicação entre sdp e sci-pan
type PanReceive struct {
	ID   int
	Time time.Time
}

// ScdPackage representa a estrutura do pacote refente ao serviço consolidador de dados
type ScdPackage struct {
	SlpOk     bool
	PanOk     bool
	GpsOk     bool
	ZoomFrame []byte
	PlateInfo jidosha.Reconhecimento
	Time      time.Time
	PanFrame  []byte
	GpsData   SglPackage
	Err       error
}

// SlpPackage representa a estrutra do pacote referente ao serviço slp
type SlpPackage struct {
	ZoomID       int
	ID           int
	ZoomFrame    *image.ImageStruct
	PlateInfo    jidosha.Reconhecimento
	JidoshaError error
}

type BufferPackage struct {
	PlateInfo jidosha.Reconhecimento
	Frame     *image.ImageStruct
}

type BufferPlate map[string]BufferPackage

// ScdBuffer representa o tipo de mapa para controle de buffer no scd
type ScdBuffer map[int]ScdPackage

// PlateBuffer representa o tipo de mapa para controle de placa no scd
type PlateBuffer map[string]time.Time

// BufferGPS representa o tipo de mapa para controle de dados do GPS
type BufferGPS map[time.Time]SglPackage

var (

	/*	Canais para comunicação entre serviços
		Nome 		Estrutura 			Origem 		Destino
		panCh		[]byte				cam			sci-pan
		panToScd	Msg					sci-pan		scd
		sdpToPan	PanRecive			sdp			sci-pan
		sdpToSlp	Msg					sdp			slp
		slpToScd	SlpPackage			slp			scd
		zoomCh		[]byte				cam			sci-zoom
		zoomToSdp	*image.ImageStruct	sci-zoom	sdp
		sglToScd	SglPackage			sgl			scd
		slpToSgl	MsgSgl			slp			sgl
	*/

	panCh     = make(chan []byte, defaults.BufferChannel)
	panToScd  = make(chan Msg, defaults.BufferChannel)
	sdpToPan  = make(chan PanReceive, defaults.BufferChannel)
	sdpToSlp  = make(chan Msg, defaults.BufferChannel)
	slpToScd  = make(chan SlpPackage, defaults.BufferChannel)
	zoomCh    = make(chan []byte, defaults.BufferChannel)
	zoomToSdp = make(chan *image.ImageZoomID, defaults.BufferChannel)
	sglToScd  = make(chan SglPackage, defaults.BufferChannel)
	slpToSgl  = make(chan MsgSgl, defaults.BufferChannel)
)

// GetChanPan retorna o canal para comunicação entre cam pan e sci-pan
func GetChanPan() chan []byte {
	return panCh
}

// GetPanToScd retorna o canal para comunicação entre scipan e scd
func GetPanToScd() chan Msg {
	return panToScd
}

// SendPanToScd envia dados para o canal entre scipan e sdc
func SendPanToScd(msg Msg) {
	panToScd <- msg
}

// GetSdpToPan retorna o canal para comunicação entre sdp e sci-pan
func GetSdpToPan() chan PanReceive {
	return sdpToPan
}

// SendSdpToPan envia dados para o canal entre sdp e sci-pan
func SendSdpToPan(data PanReceive) {
	sdpToPan <- data
}

// GetSdpToSlp retorna o canal para comunicação entre sdp e slp
func GetSdpToSlp() chan Msg {
	return sdpToSlp
}

// SendResultSlp envia dados para o canal entre sdp e slp
func SendResultSlp(msg Msg) {
	sdpToSlp <- msg
}

// GetSlpToScdCh retorna o canal para comunicação entre slp e scd
func GetSlpToScdCh() chan SlpPackage {
	return slpToScd
}

// SendSlpToScd envia dados para o canal entre slp e scd
func SendSlpToScd(msg SlpPackage) {
	slpToScd <- msg
}

// SendResultSci envia dados para o canal entre sdp e sci-pan
func SendResultSci(msg PanReceive) {
	sdpToPan <- msg
}

// GetChanZoom retorna o canal de comunicação entre camera (zoom) e sci-zoom
func GetChanZoom() chan []byte {
	return zoomCh
}

// GetZoomToSdpCh retorna o canal para comunicação entre scizoom e sdp
func GetZoomToSdpCh() chan *image.ImageZoomID {
	return zoomToSdp
}

// SendResultSdp envia dados para o canal entre scizoom e sdp
func SendResultSdp(msg *image.ImageZoomID) {
	zoomToSdp <- msg
}

// GetChanSgl retorna o canal para comunicação entre sgl e scd
func GetChanSgl() chan SglPackage {
	return sglToScd
}

// SendResultSglToScd envia dados para o canal entre sgl e scd
func SendResultSglToScd(msg SglPackage) {
	sglToScd <- msg
}

// GetChanSlp retorna o canal para comunicação entre o slp e sgl
func GetChanSlp() chan MsgSgl {
	return slpToSgl
}

// SendResultSlpToSgl envia dados de tempo para o canal entre o slp e sgl
func SendResultSlpToSgl(msg MsgSgl) {
	slpToSgl <- msg
}
