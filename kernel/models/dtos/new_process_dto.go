package dtos

type NewProcessDTO struct {
	PID  int    `json:"pid"`
	Size int    `json:"size"`
	File string `json:"file"`
}
