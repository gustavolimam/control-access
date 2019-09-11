package web

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/gustavolimam/control-access/src/components/defaults"
	"github.com/gustavolimam/control-access/src/components/log"
)

const (
	logService log.Service = "WEB"
)

// WebSys estrutura responsável por criar as variavéis utilizadas pelo objeto
type WebSys struct {
	port string
}

// New é a função que inicializa o objeto utilizado na função de start do server
func New() *WebSys {
	log.Log(logService, "Criado serviço")

	web := new(WebSys)
	web.port = ":666"

	return web
}

// Run função que inicia o front end, definindo a porta para acesso e chamando a api principal.
func (ws *WebSys) Run() {
	log.Log(logService, "Iniciado serviço")
	// Criação da variavel de rotas HTTP
	router := mux.NewRouter()
	if router == nil {
		log.Log(logService, "Falha na criação de novo roteador: Objeto vazio")
	}

	// api := router.PathPrefix("/api/").Subrouter()
	// ws.gpsAPIEndPoints(api)
	// ws.configAPIEndPoints(api)

	// Carrega os arquivos estáticos do Front
	fs := http.FileServer(http.Dir(path.Join(defaults.GetPath(), "client", "build")))
	router.PathPrefix("/").Handler(fs)
	if fs == nil {
		log.Log(logService, "Falha na geração do handler: Objeto vazio")
	}

	log.Log(logService, "Server running in port:", ws.port)
	if err := http.ListenAndServe(ws.port, router); err != nil {
		log.Fatal(logService, "Erro ao tentar iniciar servidor - ", err)
	}
}
