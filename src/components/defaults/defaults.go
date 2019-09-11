package defaults

import (
	"os"
	"time"
)

// EventoVeiculo estrutura que defini os dados que são utilizar para criar o evento de entrada de veículo
type EventoVeiculo struct {
	Placa    string
	Tempo    time.Time
	Portaria string
}

// GetPath função que retorna o diretório do sistema
func GetPath() string {
	path := os.Getenv("ControleAcessoSys")

	return path
}
