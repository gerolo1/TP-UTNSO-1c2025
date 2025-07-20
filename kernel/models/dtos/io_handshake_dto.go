package dtos

type IOHandshakeDTO struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
