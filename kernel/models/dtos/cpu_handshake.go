package dtos

type CpuHandshakeDTO struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Id   string `json:"id"`
}
