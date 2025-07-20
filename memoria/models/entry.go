package models

type Entry struct {
	Key   int // valor de la entrada, pasado por la CPU
	Value int // es una pagina o un frame (apunta al lugar en memoria real, no el del bitmap)
	ID    int // ID unico entre para un proceso, para poder identificar en SWAP
}
