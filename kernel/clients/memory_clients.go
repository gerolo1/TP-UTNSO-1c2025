package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
	"github.com/sisoputnfrba/tp-golang/kernel/models/dtos"
	"net/http"
	"net/url"
	"strconv"
)

func AskForMemory(pcb *models.Pcb) *http.Response {
	newProcess := dtos.NewProcessDTO{pcb.Pid, pcb.Size, pcb.File}
	jsonMessage, err := json.Marshal(newProcess)
	if err != nil {
		config.Logger.Error("ERROR CONVERTING TO JSON")
		panic(err)
	}

	response, err := http.Post(fmt.Sprintf("http://%s:%d/memory/processes/new", config.Config.IpMemory, config.Config.PortMemory), "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		config.Logger.Error("ERROR POSTING")
		panic(err)
	}
	return response
}

func SuspendProcess(pid int) *http.Response {
	query := url.Values{}
	query.Add("pid", strconv.Itoa(pid))

	resp, err := http.Get(fmt.Sprintf("http://%s:%d/memory/processes/suspend?%s", config.Config.IpMemory, config.Config.PortMemory, query.Encode()))
	if err != nil {
		config.Logger.Error("ERROR GETTING")
		panic(err)
	}
	return resp
}

func ResumeProcess(pcb *models.Pcb) *http.Response {
	query := url.Values{}
	query.Add("pid", strconv.Itoa(pcb.Pid))

	resp, err := http.Get(fmt.Sprintf("http://%s:%d/memory/processes/resume?%s", config.Config.IpMemory, config.Config.PortMemory, query.Encode()))
	if err != nil {
		config.Logger.Error("ERROR GETTING")
		panic(err)
	}
	return resp
}

func DumpProcess(pid int) *http.Response {
	dumpProcess := dtos.DumpRequestDTO{pid}
	jsonMessage, err := json.Marshal(dumpProcess)
	if err != nil {
		config.Logger.Error("ERROR CONVERTING TO JSON")
		panic(err)
	}

	response, err := http.Post(fmt.Sprintf("http://%s:%d/memory/dump", config.Config.IpMemory, config.Config.PortMemory), "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		config.Logger.Error("ERROR POSTING")
		panic(err)
	}
	return response
}

func ExitMemory(pid int) *http.Response {
	query := url.Values{}
	query.Add("pid", strconv.Itoa(pid))

	resp, err := http.Get(fmt.Sprintf("http://%s:%d/memory/processes/exit?%s", config.Config.IpMemory, config.Config.PortMemory, query.Encode()))
	if err != nil {
		config.Logger.Error("ERROR GETTING")
		panic(err)
	}
	return resp
}
