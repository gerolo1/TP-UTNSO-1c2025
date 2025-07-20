package dtos

type IoRequestDTO struct {
	PID      int    `json:"pid"`
	Device   string `json:"device"`
	Duration int    `json:"duration"`
	PC       int    `json:"pc"`
}
