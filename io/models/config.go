package models

type Config struct {
	PortIO          int    `json:"port_io"`
	IPIO            string `json:"ip_io"`
	IPKernel        string `json:"ip_kernel"`
	PortKernel      int    `json:"port_kernel"`
	LogLevel        string `json:"log_level"`
}