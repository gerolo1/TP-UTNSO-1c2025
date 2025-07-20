package dtos

type KernelHandshakeDTO struct {
	Id   string
	Ip   string
	Port int
}

type InitProcRequest struct {
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}