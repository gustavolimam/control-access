package web

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gustavolimam/control-access/src/components/log"
	"github.com/gustavolimam/control-access/src/services/web/context"
)

const (
	// StatusBadJSON representa o erro de JSON inválido
	StatusBadJSON int = 10000 + iota
)

// serverErrorEnvelope é um envelope dos erros retornados
// pelo servidor
type serverErrorEnvelope struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// zipWriter é um utilitário simples para escrever
// usando um compactador
type zipWriter struct {
	http.ResponseWriter
	io.Writer
}

// Permite executar handlers de HTTP após funções de middleware. Inclui a execução de
// alguns middlewares padrões. Note que os middlewares são executados recursivamente e
// portanto na ordem reversa à fornecida
func handleWith(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	h = compressHandler(h)
	h = defaultHeadersHandler(h)
	for _, m := range middleware {
		h = m(h)
	}
	h = contextHandler(h)
	return h
}

// Similar ao handleWith, porém sem compressHandler para correto funcionamento do stream de vídeo
func handleWith2(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	h = defaultHeadersHandler(h)
	for _, m := range middleware {
		h = m(h)
	}
	h = contextHandler(h)
	return h
}

func (w *zipWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write escreve um dado usando o Writer apropriado
// (supostamente um compactador)
func (w *zipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// compressHandler é um middleware responsável por compactar
// a resposta de acordo com as configurações do cliente.
func compressHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	L:
		for _, enc := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			switch strings.TrimSpace(enc) {
			case "gzip":
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Add("Vary", "Accept-Encoding")

				gw := gzip.NewWriter(w)
				defer gw.Close()

				w = &zipWriter{
					Writer:         gw,
					ResponseWriter: w,
				}
				break L
			case "deflate":
				w.Header().Set("Content-Encoding", "deflate")
				w.Header().Add("Vary", "Accept-Encoding")

				fw, _ := flate.NewWriter(w, flate.DefaultCompression)
				defer fw.Close()

				w = &zipWriter{
					Writer:         fw,
					ResponseWriter: w,
				}
				break L
			}
		}

		h.ServeHTTP(w, r)
	})
}

// contextHandler é um middleware responsável por assegurar que um contexto
// apropriado seja criado para o Request. Deveria ser sempre um dos primeiros
// middlewares a serem executados
func contextHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := context.NewEmpty()
		// configura alguns valores padrões para uso posterior
		context.Attach(c, r)
		h.ServeHTTP(w, r)
		context.Clear(r)
	})
}

// defaultHeadersHandler é um middleware que inclui headers padrão
// à resposta HTTP
func defaultHeadersHandler(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		h.ServeHTTP(w, r)
	})
}

// serveDone envia um objeto DONE para o cliente
func serveDone(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{\"status\" : \"done\"}")); err != nil {
		log.Log(logService, "Erro ao tentar escrever resultado DONE para o cliente - ", err)
	}
}

// serveNotFound envia um erro 404 para o cliente
func serveNotFound(w http.ResponseWriter, msg string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	env := serverErrorEnvelope{
		Error:   "not-found",
		Message: fmt.Sprintf(msg, args...),
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Log(logService, "Erro em json.NewEncoder (serveNotFound): ", err.Error())
	}
}

// serveNoPermissionError envia um erro forbidden para o cliente
func serveNoPermissionError(w http.ResponseWriter) {
	serveError(w, http.StatusForbidden, "você não tem permissão de acessar este recurso")
}

// serveInternalError envia um erro interno para o cliente
func serveInternalError(w http.ResponseWriter, msg string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	env := serverErrorEnvelope{
		Error:   "internal",
		Message: fmt.Sprintf(msg, args...),
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Log(logService, "Erro em json.NewEncoder (serveInternalError): ", err.Error())
	}
}

// serveBadRequest envia um erro 400 BAD para o cliente
func serveBadRequest(w http.ResponseWriter, msg string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	env := serverErrorEnvelope{
		Error:   "bad-request",
		Message: fmt.Sprintf(msg, args...),
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Log(logService, "Erro em json.NewEncoder (serveBadRequest): ", err.Error())
	}
}

// serveBadJSON envia um erro de JSON inválido (400 BAD) para o cliente
func serveBadJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	env := serverErrorEnvelope{
		Error:   "bad-json",
		Message: "o JSON não é válido",
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Log(logService, "Erro em json.NewEncoder (serveBadJSON): ", err.Error())
	}
}

// serveCustomError envia um erro do tipo identifier com código status para
// o cliente.
func serveCustomError(w http.ResponseWriter, status int, identifier string, msg string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	env := serverErrorEnvelope{
		Error:   identifier,
		Message: fmt.Sprintf(msg, args...),
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Log(logService, "Erro em json.NewEncoder (serveCustomError): ", err.Error())
	}
}

// serveError é uma função auxiliar que envia erros e seta o indentificador
// automaticamente
func serveError(w http.ResponseWriter, status int, msg string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	errorMessage := ""
	switch status {
	case http.StatusBadRequest:
		errorMessage = "bad-request"
	case http.StatusNotFound:
		errorMessage = "not-found"
	case http.StatusForbidden:
		errorMessage = "forbidden"
	case http.StatusUnauthorized:
		errorMessage = "unauthorized"
	default:
		errorMessage = "internal"
	}
	env := serverErrorEnvelope{
		Error:   errorMessage,
		Message: fmt.Sprintf(msg, args...),
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Log(logService, "Erro em json.NewEncoder (serveError): ", err.Error())
	}
}

// serveResult converte um objeto para formato JSON e o envia para o cliente
func serveResult(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Log(logService, "Erro em json.NewEncoder (serveResult): ", err.Error())
	}
}

// serveSendFile envia arquivo para o cliente
func serveSendFile(w http.ResponseWriter, path string) {
	openfile, err := os.Open(path)
	if err != nil {
		log.Log(logService, "Não foi possível abrir o arquivo "+path+": ", err.Error())
		serveInternalError(w, "Não foi possível abrir o arquivo: %v", err)
		return
	}
	defer openfile.Close() //Close after function return

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	//Copy the headers into the fileHeader buffer
	if _, err := openfile.Read(fileHeader); err != nil {
		log.Log(logService, "serveSendFile - Erro ao tentar ler fileHeader: ", err)
	}
	//Get content type of file
	fileContentType := http.DetectContentType(fileHeader)

	//Get the file size
	fileStat, err := openfile.Stat() //Get info from file
	if err != nil {
		log.Log(logService, "Não foi possível obter informações do arquivo "+path+": ", err.Error())
	}
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	if _, err := openfile.Seek(0, 0); err != nil {
		log.Log(logService, "serveSendFile - Erro ao ler arquivo via openfile.Seek: ", err.Error())
	}
	//'Copy' the file to the client
	if _, err := io.Copy(w, openfile); err != nil {
		log.Log(logService, "serveSendFile - Erro ao copiar arquivo via io.Copy: ", err.Error())
	}
}

// decodifica interpreta os dados recebidos do frontend e retorna a estrutura correspondente
// através da interface v.
func decodifica(w http.ResponseWriter, r *http.Request, v interface{}) error {

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		log.Log(logService, "Não foi possível decodificar JSON em Decodifica(): ", err.Error())
		serveBadJSON(w)
		return err
	}

	return nil
}

// salvaMudanca atualiza o arquivo config.json com as alterações de parâmetros.
// Caso o módulo metrológico não esteja selado, atualiza o MET também
// func salvaMudanca(config *config.Cfg, w http.ResponseWriter) error {
// 	if err := store.SaveSysConfig(config); err != nil {
// 		log.Log(logService, "Não foi possível salvar as configurações: ", err.Error())
// 		serveInternalError(w, "não foi possível salvar as configurações: %v", err)
// 		return err
// 	}
// 	return nil
// }
