package storage

import (
	"path"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gustavolimam/control-access/src/components/defaults"
	"github.com/gustavolimam/control-access/src/components/log"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const (
	logService log.Service = "FIRESTORE"
)

// RegistroVeicular estrutura à ser enviado para o BD
type RegistroVeicular struct {
	Placa    string    `json:"placa,omitempty"`
	Tempo    time.Time `json:"time,omitempty"`
	Portaria string    `json:"portaria,omitempty"`
}

// SendEntryToDB função responsável por criar conexão com o banco e enviar os dados de Evento de Entrada.
func SendEntryToDB(event defaults.EventoVeiculo) error {
	log.Log(logService, "Enviando registro de entrada de veiculo para base Firestore")

	// Retorna a data atual
	tempo := time.Now().Format("-2006-01-02")

	// Inicializando o Firebase
	opt := option.WithCredentialsFile(path.Join("util", "controle-acesso-port-firebase-adminsdk-ts97s-dec0edb44a.json"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	// Criando a conexão com o banco Firestore
	client, err := app.Firestore(context.Background())
	if err != nil {
		client.Close()
		return err
	}

	// Enviando informação para o Database - Firestore
	_, _, err = client.Collection("registro-veiculos"+tempo).Add(context.Background(), &RegistroVeicular{Placa: event.Placa,
		Tempo: event.Tempo, Portaria: event.Portaria})
	if err != nil {
		log.Log(logService, "Erro ao tentar enviar registro de entrada de veiculo para o firestore - ", err)
		client.Close()
		return err
	}

	return nil
}

//SendExitToDB função responsável por enviar os dados de evento de saída para o banco de dados
func SendExitToDB(event defaults.EventoVeiculo) error {
	log.Log(logService, "Enviando registro de saída de veiculo para base Firestore")

	// Retorna a data atual
	tempo := time.Now().Format("-2006-01-02")

	// Inicializando o Firebase
	opt := option.WithCredentialsFile(path.Join("util", "controle-acesso-port-firebase-adminsdk-ts97s-dec0edb44a.json"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Log(logService, "Erro inicializar app:", err)
		return err
	}

	// Criando a conexão com o banco Firestore
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Log(logService, "Erro ao inicializar o Firestore: ", err)
		client.Close()
		return err
	}

	Document, err := client.Collection("registro-veiculos"+tempo).Where("Placa", "==", "EXI7254").Documents(context.Background()).GetAll()
	if err != nil {
		client.Close()
		return err
	}

	_, err = Document[0].Ref.Set(context.Background(), map[string]interface{}{"TempoSaida": event.Tempo, "PortariaSaida": event.Portaria}, firestore.MergeAll)
	if err != nil {
		client.Close()
		return err
	}

	return nil
}
