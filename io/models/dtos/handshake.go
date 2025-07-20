package dtos

// Petición de IO del Kernel al IO
// request que el Kernel le envía al módulo IO para op I/O
type IORequestDTO struct {
	PID      int `json:"pid"`
	Duration int `json:"duration"` // milisegundos
}

//Confirm handshake
type ResponseDTO struct {
	Message string `json:"message"`
}

// Handshake inicial con Kernel
type IOHandshakeDTO struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
