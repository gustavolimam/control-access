package main

import (
	"fmt"

	"github.com/gustavolimam/control-access/src/components/config"
	"github.com/gustavolimam/control-access/src/components/log"
	"github.com/gustavolimam/control-access/src/components/report"
	"github.com/gustavolimam/control-access/src/services/events"
	"github.com/gustavolimam/control-access/src/services/web"
)

const (
	logService log.Service = "MAIN"
)

func main() {
	fmt.Println("Iniciando sistema de controle de acesso FACENS")

	// Carrega o arquivo config.json
	if err := config.SetupConfig(); err != nil {
		log.Fatal(logService, "Erro ao carregar configurações: ", err)
	}

	// Cria o arquivo de log
	_, err := log.CreateLogFile()
	if err != nil {
		log.Fatal(logService, "Erro ao criar arquivo de log: ", err)
	}

	// Start events service
	if ev := events.New(); ev == nil {
		log.Fatal(logService, "Erro ao criar Serviço de Eventos")
	} else {
		go ev.Run()
	}

	// Start web service
	if ws := web.New(); ws == nil {
		log.Fatal(logService, "Erro ao criar Sistema Web")
	} else {
		go ws.Run()
	}

	// Escuta o canal de error para uma possivel reinicialização
	chError := report.GetErrorReportCh()

	for {
		err := <-chError
		log.Fatal(logService, err)
	}
}
