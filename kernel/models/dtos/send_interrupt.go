package dtos

type SendInterrupt struct {
	PID    int    `json:"pid"`
	Motivo string `json:"motivo"`
}
