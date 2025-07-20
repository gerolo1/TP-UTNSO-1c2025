package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/models"
)

type VarConfig struct {
	SchedulerAlgorithm    string  `json:"scheduler_algorithm"`
	ReadyIngressAlgorithm string  `json:"ready_ingress_algorithm"`
	Alpha                 float64 `json:"alpha"`
	InitialEstimate       float64 `json:"initial_estimate"`
	SuspensionTime        int     `json:"suspension_time"`
}

type FinalConfig struct {
	IpMemory   string `json:"ip_memory"`
	PortMemory int    `json:"port_memory"`
	LogLevel   string `json:"log_level"`
	PortKernel int    `json:"port_kernel"`
	IpKernel   string `json:"ip_kernel"`
}

var Config *models.Config

func InitConfiguration() string {
	varConfig := &VarConfig{}
	finalConfig := &FinalConfig{}
	configType := 0

	fmt.Println("1 - PLANIFICACION CORTO PLAZO - CASO FIFO")
	fmt.Println("2 - PLANIFICACION CORTO PLAZO - CASO SJF")
	fmt.Println("3 - PLANIFICACION CORTO PLAZO - CASO SRT")
	fmt.Println("4 - PLANIFICACION MEDIANO/LARGO PLAZO - CASO FIFO")
	fmt.Println("5 - PLANIFICACION MEDIANO/LARGO PLAZO - CASO PMCP")
	fmt.Println("6 - MEMORIA SWAP")
	fmt.Println("7 - MEMORIA CACHE")
	fmt.Println("8 - MEMORIA TLB")
	fmt.Println("9 - ESTABILIDAD GENERAL")

	fmt.Scan(&configType)

	configFile, e := os.Open("kernel/config/config.json")
	if e != nil {
		panic(e)
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	e = jsonParser.Decode(finalConfig)
	if e != nil {
		panic(e)
	}

	configVar, err := os.Open("kernel/config/" + strconv.Itoa(configType) + ".json")
	if err != nil {
		panic(err)
	}

	defer configVar.Close()

	jsonParser = json.NewDecoder(configVar)
	e = jsonParser.Decode(varConfig)
	if e != nil {
		panic(e)
	}

	Config = &models.Config{
		IpMemory:              finalConfig.IpMemory,
		PortMemory:            finalConfig.PortMemory,
		IpKernel:              finalConfig.IpKernel,
		PortKernel:            finalConfig.PortKernel,
		LogLevel:              finalConfig.LogLevel,
		SchedulerAlgorithm:    varConfig.SchedulerAlgorithm,
		ReadyIngressAlgorithm: varConfig.ReadyIngressAlgorithm,
		Alpha:                 varConfig.Alpha,
		InitialEstimate:       varConfig.InitialEstimate,
		SuspensionTime:        varConfig.SuspensionTime,
	}

	switch configType {
	case 1, 2, 3:
		return "PLANI_CORTO_PLAZO"
	case 4, 5:
		return "PLANI_LYM_PLAZO"
	case 6:
		return "MEMORIA_IO"
	case 7:
		return "MEMORIA_BASE"
	case 8:
		return "MEMORIA_BASE_TLB"
	default:
		return "ESTABILIDAD_GENERAL"
	}
}

func FirstProcessSize(nombre string) int {
	switch nombre {
	case "PLANI_CORTO_PLAZO", "PLANI_LYM_PLAZO", "ESTABILIDAD_GENERAL":
		return 0
	case "MEMORIA_IO":
		return 90
	default:
		return 256
	}
}
