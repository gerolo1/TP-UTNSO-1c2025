package dtos

type WriterDTO struct {
	PID   int    `json:"pid"`
	DF    int    `json:"df"`
	Value []byte `json:"value"`
}
