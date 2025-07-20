package dtos

type WriterDTO struct {
	PID   int    `json:"pid"`
	DF    int    `json:"df"`    // dirección física
	Value []byte `json:"value"` // valor a escribir
}
