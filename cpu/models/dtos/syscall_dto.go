package dtos

//para IO
type IORequest struct {
	PID      int    `json:"pid"`
	Device   string `json:"device"`
	Duration int    `json:"duration"`
	PC       int    `json:"pc"`
}


//para exit
type ExitProcess struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}