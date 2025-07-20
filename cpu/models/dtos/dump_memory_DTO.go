package dtos

type DumpRequestDTO struct {
	PID int `json:"pid"`
	Pc  int `json:"pc"` 
}

type DumpResponseDTO struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}
