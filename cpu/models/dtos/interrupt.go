package dtos

type Interrupt struct {
	PID    int    `json:"pid"`
	Motivo string `json:"motivo"` 
}

type ProcessInterrupt struct {
	PID    int    `json:"pid"`
	PC     int    `json:"pc"`
	Motivo string `json:"motivo"`
}