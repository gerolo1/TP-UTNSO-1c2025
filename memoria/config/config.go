package config

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/models"
	"os"
	"strconv"
)

type VarConfig struct {
	MemorySize     int `json:"memory_size"`
	PageSize       int `json:"page_size"`
	EntriesPerPage int `json:"entries_per_page"`
	NumberOfLevels int `json:"number_of_levels"`
	MemoryDelay    int `json:"memory_delay"`
	SwapDelay      int `json:"swap_delay"`
}

type FinalConfig struct {
	PortMemory   int    `json:"port_memory"`
	IpMemory     string `json:"ip_memory"`
	SwapfilePath string `json:"swapfile_path"`
	LogLevel     string `json:"log_level"`
	DumpPath     string `json:"dump_path"`
	ScriptsPath  string `json:"scripts_path"`
}

func InitConfiguration() *models.Config {
	varConfig := &VarConfig{}
	finalConfig := &FinalConfig{}
	configType := 0

	fmt.Println("1 - PLANIFICACION CORTO PLAZO")
	fmt.Println("2 - PLANIFICACION LARGO PLAZO")
	fmt.Println("3 - MEMORIA SWAP")
	fmt.Println("4 - MEMORIA CACHE")
	fmt.Println("5 - MEMORIA TLB")
	fmt.Println("6 - ESTABILIDAD GENERAL")
	fmt.Scan(&configType)

	configFile, err := os.Open("memoria/config/config.json")
	if err != nil {
		panic(err)
	}

	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			panic(err)
		}
	}(configFile)

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(finalConfig)
	if err != nil {
		panic(err)
	}

	configVar, err := os.Open("memoria/config/json/" + strconv.Itoa(configType) + ".json")
	if err != nil {
		panic(err)
	}

	defer func(configVar *os.File) {
		e := configVar.Close()
		if e != nil {
			panic(e)
		}
	}(configVar)

	jsonParser = json.NewDecoder(configVar)
	err = jsonParser.Decode(varConfig)
	if err != nil {
		panic(err)
	}

	config := &models.Config{
		MemorySize:     varConfig.MemorySize,
		PageSize:       varConfig.PageSize,
		EntriesPerPage: varConfig.EntriesPerPage,
		NumberOfLevels: varConfig.NumberOfLevels,
		MemoryDelay:    varConfig.MemoryDelay,
		SwapDelay:      varConfig.SwapDelay,
		PortMemory:     finalConfig.PortMemory,
		IpMemory:       finalConfig.IpMemory,
		SwapfilePath:   finalConfig.SwapfilePath,
		LogLevel:       finalConfig.LogLevel,
		DumpPath:       finalConfig.DumpPath,
		ScriptsPath:    finalConfig.ScriptsPath,
	}

	return config
}
