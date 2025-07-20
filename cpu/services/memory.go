package services

import (
	"fmt"
	"math"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/config"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/models/dtos"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
)

func GetInstructionFromMemory(pid int, pc int, targetIP string, targetPort int) (string, error) {
	endpoint := fmt.Sprintf("memory/instruction?pid=%d&pc=%d", pid, pc)
	var data map[string]string

	config.Logger.Debug("Obteniendo instrucción de memoria para PID " + strconv.Itoa(pid) + " y PC " + strconv.Itoa(pc))

	err := utils.GetFromModule(targetIP, targetPort, endpoint, &data)
	if err != nil {
		config.Logger.Error("Error al obtener instrucción de memoria: " + err.Error())
		return "", err
	}

	instruction, ok := data["instruction"]
	if !ok {
		return "", fmt.Errorf("campo 'instruction' no encontrado en la respuesta")
	}

	return instruction, nil
}

func GetFrameFromMemory(pid int, dl int, ip string, port int) (int, error) {
	nroPagina := dl / globals.PageSize
	indices := make([]int, globals.Levels)
	// Calcular los índices para cada nivel de la tabla de páginas

	for X := 1; X <= globals.Levels; X++ {
		divisor := int(math.Pow(float64(globals.EntriesPerTable), float64(globals.Levels-X)))
		indices[X-1] = (nroPagina / divisor) % globals.EntriesPerTable
	}

	requestBody := map[string]interface{}{
		"pid":     pid,
		"entries": indices,
	}

	var response map[string]int
	err := utils.SendToModule(ip, port, "memory/frame", requestBody, &response)
	if err != nil {
		return -1, err
	}

	frame, ok := response["frame"]
	if !ok {
		return -1, fmt.Errorf("campo 'frame' no encontrado en la respuesta")
	}

	return frame, nil
}

func ReadMemory(pid int, df int, sizeToRead int, targetIP string, targetPort int) ([]byte, error) {
	endpoint := fmt.Sprintf("memory/read?pid=%d&address=%d&size=%d", pid, df, sizeToRead)
	config.Logger.Debug("Leyendo memoria para PID " + strconv.Itoa(pid) + ", dirección física " + strconv.Itoa(df) + ", tamaño " + strconv.Itoa(sizeToRead))

	dto := dtos.ReaderDTO{
		PID:  pid,
		Size: sizeToRead,
		DF:   df,
	}

	var data struct {
		Data []byte `json:"data"`
	}
	err := utils.SendToModule(targetIP, targetPort, endpoint, dto, &data)
	if err != nil {
		config.Logger.Error("Error al leer memoria: " + err.Error())
		return nil, err
	}

	return data.Data, nil
}

func GetPageTableConfig(ip string, port int) (pageSize int, levels int, entriesPerTable int, err error) {
	endpoint := "memory/config"

	var data map[string]int
	err = utils.GetFromModule(ip, port, endpoint, &data)
	if err != nil {
		return
	}

	var ok bool
	pageSize, ok = data["page_size"]
	if !ok {
		err = fmt.Errorf("campo 'page_size' no encontrado en la respuesta")
		return
	}

	levels, ok = data["number_of_levels"]
	if !ok {
		err = fmt.Errorf("campo 'number_of_levels' no encontrado en la respuesta")
		return
	}

	entriesPerTable, ok = data["entries_per_page"]
	if !ok {
		err = fmt.Errorf("campo 'entries_per_page' no encontrado en la respuesta")
		return
	}

	return pageSize, levels, entriesPerTable, nil
}

func WriteMemory(dto dtos.WriterDTO, targetIP string, targetPort int) error {
	endpoint := "memory/write"
	config.Logger.Debug("Escribiendo en memoria para PID " + strconv.Itoa(dto.PID) +
		", DF: " + strconv.Itoa(dto.DF) + ", valor: " + string(dto.Value))

	err := utils.SendToModule(targetIP, targetPort, endpoint, dto, nil)
	if err != nil {
		config.Logger.Error("Error al escribir en memoria: " + err.Error())
		return err
	}

	return nil
}
