package buffer

import (
	"errors"
	"sync"
	"time"
)

const (
	bufferTimeoutLoop = time.Minute // tempo para verificação de dados obsoletos
	bufferTimeout     = 60          // tempo para descartar dados obsoletos - Em minutos
)

var errFrameNotFound = errors.New("Nao foi possivel obter frame")

// FrameBuffer define a estrutura do buffer circular de vídeo
type FrameBuffer struct {
	s           int
	d           []*image.ImageStruct
	w           int // índice para próxima escrita. o último elemento é w-1.
	r           int // índice do primeiro elemento.
	bufferMutex sync.Mutex
}

// NewBuffer cria um novo buffer
func NewBuffer(size int) *FrameBuffer {
	b := new(FrameBuffer)
	b.s = size
	b.d = make([]*image.ImageStruct, b.s+1)
	b.bufferMutex = sync.Mutex{}
	return b
}

// Add adiciona um elemento no buffer circular
func (b *FrameBuffer) Add(ml *image.ImageStruct) {
	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()
	b.d[b.w] = ml
	b.w++
	if b.w >= b.s+1 {
		b.w = 0
	}
	if b.r == b.w {
		b.r++
		if b.r >= b.s+1 {
			b.r = 0
		}
	}
}

// Frame retorna o frame mais próximo ao tempo t
func (b *FrameBuffer) Frame(t time.Time) (*image.ImageStruct, error) {
	var i int
	var err error
	if i, err = b.Find(t); err != nil {
		return nil, err
	}
	return b.d[i], nil
}

// Len retorna o número de elementos existentes no buffer
func (b *FrameBuffer) Len() int {
	if b.w < b.r {
		return b.s + 1 + b.w - b.r
	}
	return b.w - b.r
}

// Find busca o índice do elemento que tem o timestamp mais próximo ao tempo t
func (b *FrameBuffer) Find(t time.Time) (int, error) {
	b.bufferMutex.Lock()
	defer b.bufferMutex.Unlock()
	pFim := b.w - 1
	if pFim < 0 {
		pFim = b.s
	}

	len := b.Len()
	if len == 0 {
		return -1, errFrameNotFound
	}
	periodoTotal := b.d[pFim].Time.Sub(b.d[b.r].Time)
	if periodoTotal == 0 { // tratamento para divizão por zero
		return b.r, nil
	}
	// faz uma estimativa da posição do frame buscado
	i := int(time.Duration(len-1) * t.Sub(b.d[b.r].Time) / periodoTotal)
	if i < 0 {
		i = 0
	}
	if i > len-1 {
		i = len - 1
	}
	i += b.r
	if i >= b.s+1 {
		i -= (b.s + 1)
	}

	if b.d[i].Time.Before(t) {
		// o frame da posição estimada é anterior ao buscado
		temp := i
		for {
			// acessa o elemento seguinte
			i++
			if i >= b.s+1 {
				i = 0
			}

			// se não existe mais elemento
			if i == b.w {
				return temp, nil
			}

			// verifica o elemento
			if b.d[i].Time.Before(t) {
				// ainda é anterior ao buscado. continua
				temp = i
			} else {
				// é posterior.
				// compara o elemento imediatamente posterior e o imediatamente
				// anterior e retorna o que índice do que estiver mais próximo.
				if b.d[i].Time.Sub(t) > t.Sub(b.d[temp].Time) {
					return temp, nil
				}
				return i, nil
			}
		}
	} else {
		// o frame da posição estimada é posterior ao buscado
		temp := i
		for {
			// verifica se é o primeiro elemento
			if i == b.r {
				return temp, nil
			}

			// acessa elemento anterior
			if i == 0 {
				i = b.s
			} else {
				i--
			}

			// verifica o elemento
			if b.d[i].Time.After(t) {
				// ainda é posterior ao buscado. continua
				temp = i
			} else {
				// é anterior ao buscado
				// compara o elemento imediatamente posterior e o imediatamente
				// anterior e retorna o que índice do que estiver mais próximo.
				if b.d[temp].Time.Sub(t) > t.Sub(b.d[i].Time) {
					return i, nil
				}
				return temp, nil
			}
		}
	}
}

// Frames retorna os frames no intervalo e taxa de frame especificado.
// Antes de chamar esta função deve-se garantir que o buffer já contém
// os frames no range especificado.
// func (b *FrameBuffer) Frames(inicio, fim time.Time, fps int) []*image.ImageStruct {
// 	i := b.Find(inicio)
// 	passo := time.Second / time.Duration(fps)
// 	totalFrame := int(fim.Sub(inicio)*time.Duration(fps)/time.Second + 1)
// 	retorno := make([]*image.ImageStruct, totalFrame)
// 	tempo := inicio
// 	x := 0
// 	for tempo.Before(b.d[i].Time) {
// 		retorno[x] = b.d[i]
// 		x++
// 		if x >= totalFrame {
// 			return retorno
// 		}
// 		tempo = tempo.Add(passo)
// 	}
// 	iAnt := 0
// 	for ; x < totalFrame; x++ {
// 		for b.d[i].Time.Before(tempo) {
// 			iAnt = i
// 			i++
// 			if i == b.w {
// 				for ; x < totalFrame; x++ {
// 					retorno[x] = b.d[iAnt]
// 				}
// 				return retorno
// 			}
// 			if i >= b.s+1 {
// 				i = 0
// 			}
// 		}
// 		if tempo.Sub(b.d[iAnt].Time) < b.d[i].Time.Sub(tempo) {
// 			retorno[x] = b.d[iAnt]
// 		} else {
// 			retorno[x] = b.d[i]
// 		}
// 		tempo = tempo.Add(passo)
// 	}
// 	return retorno
// }
