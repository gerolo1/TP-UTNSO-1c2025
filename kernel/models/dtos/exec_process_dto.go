package dtos

type ExecProcessDTO struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}
