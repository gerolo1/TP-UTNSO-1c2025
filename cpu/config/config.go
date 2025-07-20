package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/models"
)

var Config *models.Config

type FixedConfig struct {
	PortCPU    int    `json:"port_cpu"`
	IPCPU      string `json:"ip_cpu"`
	IPMemory   string `json:"ip_memory"`
	PortMemory int    `json:"port_memory"`
	IPKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	LogLevel   string `json:"log_level"`
}

type VariableConfig struct {
	TLBEntries       int    `json:"tlb_entries"`
	TLBReplacement   string `json:"tlb_replacement"`
	CacheEntries     int    `json:"cache_entries"`
	CacheReplacement string `json:"cache_replacement"`
	CacheDelay       int    `json:"cache_delay"`
}

func InitConfiguration() *models.Config {
	fmt.Println("1 - PLANIFICACION CORTO/MEDIANO/LARGO PLAZO")
	fmt.Println("2 - MEMORIA SWAP")
	fmt.Println("3 - MEMORIA CACHE - CLOCK")
	fmt.Println("4 - MEMORIA CACHE - CLOCK-M")
	fmt.Println("5 - MEMORIA TLB - CASO FIFO")
	fmt.Println("6 - MEMORIA TLB - CASO LRU")
	fmt.Println("7 - ESTABILIDAD GENERAL - CPU 1")
	fmt.Println("8 - ESTABILIDAD GENERAL - CPU 2")
	fmt.Println("9 - ESTABILIDAD GENERAL - CPU 3")
	fmt.Println("10 - ESTABILIDAD GENERAL - CPU 4")
	var configName string
	fmt.Scan(&configName)
	configName += ".json"

	// Abrir archivo fijo
	configFile, err := os.Open("cpu/config/config.json")
	if err != nil {
		log.Fatal("Error abriendo config.json:", err)
	}
	defer configFile.Close()

	// Abrir archivo variable
	configFile2, err := os.Open("cpu/config/" + configName)
	if err != nil {
		log.Fatal("Error abriendo archivo variable:", err)
	}
	defer configFile2.Close()

	// Decodificar archivos en structs separados
	var fixed FixedConfig
	var variable VariableConfig

	if err := json.NewDecoder(configFile).Decode(&fixed); err != nil {
		log.Fatal("Error al decodificar config.json:", err)
	}
	if err := json.NewDecoder(configFile2).Decode(&variable); err != nil {
		log.Fatal("Error al decodificar archivo variable:", err)
	}

	// Unir ambas configuraciones
	finalConfig := &models.Config{
		PortCPU:          fixed.PortCPU,
		IPCPU:            fixed.IPCPU,
		IPMemory:         fixed.IPMemory,
		PortMemory:       fixed.PortMemory,
		IPKernel:         fixed.IPKernel,
		PortKernel:       fixed.PortKernel,
		LogLevel:         fixed.LogLevel,
		TLBEntries:       variable.TLBEntries,
		TLBReplacement:   variable.TLBReplacement,
		CacheEntries:     variable.CacheEntries,
		CacheReplacement: variable.CacheReplacement,
		CacheDelay:       variable.CacheDelay,
	}

	return finalConfig
}

func SaveConfiguration(config *models.Config) {
	configFile, err := os.Create("cpu/config/config.json")
	if err != nil {
		log.Fatal("No se pudo guardar la configuración:", err)
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ") // Esto es solo para que el JSON quede legible
	err = encoder.Encode(config)
	if err != nil {
		log.Fatal("Error al escribir configuración:", err)
	}
}
