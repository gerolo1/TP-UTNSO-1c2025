package models

type Config struct {
	IpMemory              string  `json:"ip_memory"`
	PortMemory            int     `json:"port_memory"`
	SchedulerAlgorithm    string  `json:"scheduler_algorithm"`
	ReadyIngressAlgorithm string  `json:"ready_ingress_algorithm"`
	Alpha                 float64 `json:"alpha"`
	InitialEstimate       float64 `json:"initial_estimate"`
	SuspensionTime        int     `json:"suspension_time"`
	LogLevel              string  `json:"log_level"`
	PortKernel            int     `json:"port_kernel"`
	IpKernel              string  `json:"ip_kernel"`
}
