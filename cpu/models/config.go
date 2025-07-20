package models

type Config struct {
	PortCPU         int    `json:"port_cpu"`
	IPCPU           string `json:"ip_cpu"`
	IPMemory        string `json:"ip_memory"`
	PortMemory      int    `json:"port_memory"`
	IPKernel        string `json:"ip_kernel"`
	PortKernel      int    `json:"port_kernel"`
	TLBEntries      int    `json:"tlb_entries"`
	TLBReplacement  string `json:"tlb_replacement"`
	CacheEntries    int    `json:"cache_entries"`
	CacheReplacement string `json:"cache_replacement"`
	CacheDelay      int    `json:"cache_delay"`
	LogLevel        string `json:"log_level"`
}
