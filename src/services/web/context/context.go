// O pacote context implementa o utilitário de contexto para compartilhamento
// de informações entre APIs que estejam lidando com uma requisição HTTP única
package context

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

type key int

var (
	ctxsMu sync.Mutex
	ctxs   = make(map[*http.Request]context.Context)

	DbNameKey key = 0 // chave para acessar o nome do banco de dados
	UserKey   key = 1 // chave para acessar o usuário conectado
)

// New cria um contexto associado à uma requisição. Se um contexto
// foi previamente criado e associado à mesma requisição, este
// será usado
func New(r *http.Request) context.Context {
	ctxsMu.Lock()
	c := ctxs[r]
	ctxsMu.Unlock()

	if c == nil {
		// no context attached, so create a new one
		c = NewEmpty()
		Attach(c, r)
	}
	return c
}

// NewEmpty cria um contexto vazio e não associado a uma
// requisição
func NewEmpty() context.Context {
	return context.TODO()
}

// Attach associa um contexto a uma requisição. Note que se
// já houver um contexto associado à mesma requisição, este
// será sobrescrito
func Attach(ctx context.Context, r *http.Request) {
	ctxsMu.Lock()
	defer ctxsMu.Unlock()
	ctxs[r] = ctx
}

// Clear apaga o contexto associado à requisição
func Clear(r *http.Request) {
	ctxsMu.Lock()
	defer ctxsMu.Unlock()
	delete(ctxs, r)
}
