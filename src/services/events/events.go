package events

import (
	"time"

	"github.com/gustavolimam/control-access/src/components/defaults"
	"github.com/gustavolimam/control-access/src/components/log"
	"github.com/gustavolimam/control-access/src/components/storage"
)

const (
	logService log.Service = "EVENTS"
)

// EventSys estrutura do serviço de eventos
type EventSys struct {
	WebCh chan defaults.EventoVeiculo
}

// New instancia o serviço de eventos
func New() *EventSys {
	log.Log(logService, "Criado Serviço")

	ev := new(EventSys)
	ev.WebCh = make(chan defaults.EventoVeiculo, 300)

	return ev
}

// Run função responsável por ficar enviando os dados de evento de entrada e saída da portaria.
func (ev *EventSys) Run() {
	log.Log(logService, "Iniando a recepção de dados de entrada ou saida da portaria")

	// Simula o envio de dados a cada 2 minutos
	tickerEntrada := time.NewTicker(20 * time.Second)
	go func() {
		for _ = range tickerEntrada.C {
			log.Log(logService, "Novo evento de entrada de veiculo")
			// Criação de dados fakes para simular a entrada de um veiculo
			dadosTest := defaults.EventoVeiculo{
				Placa:    "EXI7254",
				Tempo:    time.Now(),
				Portaria: "Principal",
			}

			if err := storage.SendEntryToDB(dadosTest); err != nil {
				log.Log(logService, "Erro ao tentar enviar informação de entrada - erro: ", err)
			}
		}
	}()

	// Simula o envio de dados a cada 3 minutos
	tickerSaida := time.NewTicker(30 * time.Second)
	go func() {
		for _ = range tickerSaida.C {
			log.Log(logService, "Novo evento de saída de veiculo")
			// Criação de dados fakes para simular a entrada de um veiculo
			dadosTest := defaults.EventoVeiculo{
				Placa:    "EXI7254",
				Tempo:    time.Now(),
				Portaria: "Iguatemi",
			}

			if err := storage.SendExitToDB(dadosTest); err != nil {
				log.Log(logService, "Erro ao tentar enviar informação de saída - erro: ", err)
			}
		}
	}()

	// for {
	// 	newEvent := <-dadosFakes
	// 	log.Log(logService, "Nova entrada de veículo detecta na portaria principal")

	// 	// Envia as informações de evento de entrada para serem mostrada no Front
	// 	ev.WebCh <- newEvent
	// 	// Envia as informações de evento de entrada para serem salvas no banco de dados
	// 	storage.SendEntryToDB(newEvent)
	// }
}
