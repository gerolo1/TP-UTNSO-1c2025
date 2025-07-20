package dtos

type Handshake struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Id     string `json:"id"`
}
