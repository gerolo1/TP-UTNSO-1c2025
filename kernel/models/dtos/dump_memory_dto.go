package dtos

type DumpRequestDTO struct {
	PID int `json:"pid"`
}

type DumpReceiveDTO struct {
	PID int `json:"pid"`
	Pc  int `json:"pc"`
}
