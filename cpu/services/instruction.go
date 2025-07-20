package services

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/models/dtos"
	"strconv"
	"strings"

	"github.com/sisoputnfrba/tp-golang/cpu/config"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func ProcessInstruction(pid int, pc int) {

	//3.b Chequear si es un proceso nuevo, y desalojar TLB y Cache
	if globals.PID != pid {
		//Flush de lo q tenia a memoria
		FlushModifiedPagesToMemory(pid)
		// Cambió el proceso
		globals.Tlb.Entries = nil
		globals.Tlb.OrderCont = 0
		globals.Tlb.AccesCont = 0
		globals.Cache.Entries = nil
		globals.Cache.ClockPointer = 0

		globals.PC = pc
		globals.PID = pid
	}

	//3.c Obtener instrucción en memoria en base al PC y PID
	instruction, errNextInstruction := GetInstructionFromMemory(pid, pc, globals.ConfigFile.IPMemory, globals.ConfigFile.PortMemory)
	if errNextInstruction != nil {
		config.Logger.Debug("Error al obtener la instrucción de memoria")
		return
	} else {
		//config.Logger.Debug("Siguiente instrucción: " + instruction)
		config.Logger.Info("## PID: " + strconv.Itoa(pid) + " - Fetch - Program Counter: " + strconv.Itoa(pc))
	}

	/* 	instruction := map[string]string{
		"instruction": "WRITE 101212 EJEMPLO_DE_ENUNCIADO",
	} */

	//3.d Decode instruction
	instruction_string := instruction

	loggerString := "PID: " + strconv.Itoa(pid) + " - Ejecutando: " + instruction_string
	config.Logger.Info(loggerString)
	// Decode de READ/WRITE

	fields := strings.Fields(instruction_string)
	if len(fields) == 0 {
		return
	}

	instruction = strings.ToUpper(fields[0])

	// Decode READ/WRITE
	if instruction == "READ" || instruction == "WRITE" {

		dl := instruction_string[strings.Index(instruction_string, " ")+1 : strings.Index(instruction_string[strings.Index(instruction_string, " ")+1:], " ")+strings.Index(instruction_string, " ")+1]
		dlInt, err := strconv.Atoi(dl)

		if err != nil {
			config.Logger.Debug("Error al convertir dl a entero")
		}
		config.Logger.Debug("DlInt: " + strconv.Itoa(dlInt))

		numberPage := dlInt / globals.PageSize
		offset := dlInt % globals.PageSize

		frame := 0
		isInTLB := false
		isInCache := false
		var pageContent string

		// Si la caché está habilitada, esta la pag?
		if globals.ConfigFile.CacheEntries > 0 {
			pageContent, isInCache, frame = GetFromCache(numberPage)
			if isInCache {
				config.Logger.Info("PID: " + strconv.Itoa(pid) + " CACHE HIT - Página: " + strconv.Itoa(numberPage))
			} else {
				config.Logger.Info("PID: " + strconv.Itoa(pid) + " CACHE MISS - Página: " + strconv.Itoa(numberPage))
			}
		}

		if !isInCache {

			//Si la TLB está habilitada
			if globals.ConfigFile.TLBEntries > 0 {
				frame, isInTLB = GetFrameFromLocalTLB(numberPage)
				if !isInTLB {
					loggerString := "PID: " + strconv.Itoa(pid) + " - TLB MISS - Página: " + strconv.Itoa(numberPage)
					config.Logger.Info(loggerString)
				}
			}

			if !isInTLB {
				// Obtengo el frame de memoria teniendo en cuenta multinivel
				var errFrame error
				frame, errFrame = GetFrameFromMemory(pid, dlInt, globals.ConfigFile.IPMemory, globals.ConfigFile.PortMemory)
				if errFrame != nil {
					config.Logger.Debug("Error al obtener el frame de memoria")
				} else {
					loggerString := "PID: " + strconv.Itoa(pid) + " - OBTENER MARCO - Página: " + strconv.Itoa(numberPage) + " - Marco: " + strconv.Itoa(frame)
					config.Logger.Info(loggerString)

					if globals.ConfigFile.TLBEntries > 0 {
						// Agregar a la TLB
						AddEntry(numberPage, frame, globals.ConfigFile.TLBReplacement)
						config.Logger.Debug("Entrada agregada a la TLB: " + strconv.Itoa(numberPage) + " -> " + strconv.Itoa(frame))
					}

					// Agregar a cache
					// Si la cache está habilitada, leo desde memory la pag para agregarla
					if globals.ConfigFile.CacheEntries > 0 {
						pageContentBytes, errRead := ReadMemory(pid, frame, globals.PageSize, globals.ConfigFile.IPMemory, globals.ConfigFile.PortMemory)
						if errRead != nil {
							config.Logger.Error("Error al leer página desde memoria para cache: " + errRead.Error())
						} else {
							pageContent = string(pageContentBytes)
							AddToCache(numberPage, pageContent, false, frame+offset) // false: no fue modificada
							config.Logger.Info("PID: " + strconv.Itoa(pid) + "CACHE ADD - Página: " + strconv.Itoa(numberPage))
							config.Logger.Debug("Página cacheada tras lectura: " + strconv.Itoa(numberPage))
						}
					}
				}
			} else {
				loggerString := "PID: " + strconv.Itoa(pid) + " - TLB HIT - Página: " + strconv.Itoa(numberPage) + " - Marco: " + strconv.Itoa(frame)
				config.Logger.Info(loggerString)
			}

		}

		direccionFisica := frame + offset

		if len(fields) > 0 && strings.ToUpper(fields[0]) == "WRITE" {
			parts := strings.Split(instruction_string, " ")
			valor := parts[2]
			//Si esta habilitada cache, escribo en cache
			if globals.ConfigFile.CacheEntries > 0 {
				writeCache(numberPage, valor)
				loggerString := "PID: " + strconv.Itoa(pid) + " - Acción: ESCRIBIR - CACHE: " + strconv.Itoa(direccionFisica) + " - Valor: " + valor
				config.Logger.Info(loggerString)
			} else {
				// Escribir en memoria
				dto := dtos.WriterDTO{
					PID:   pid,
					DF:    direccionFisica,
					Value: []byte(valor),
				}

				errWrite := WriteMemory(dto, globals.ConfigFile.IPMemory, globals.ConfigFile.PortMemory)
				if errWrite != nil {
					config.Logger.Error("Error al escribir en memoria: " + errWrite.Error())
				} else {
					config.Logger.Info("PID: " + strconv.Itoa(pid) + " - Acción: ESCRIBIR - Dirección Física:" + strconv.Itoa(direccionFisica) + " - Valor: " + valor)
				}
			}
			// Agregar a cache (modificado = true)
			/*if globals.ConfigFile.CacheEntries > 0 {
				AddToCache(numberPage, valor, true, frame)
				config.Logger.Info("PID: " + strconv.Itoa(pid) + "CACHE ADD - Página: " + strconv.Itoa(numberPage))
				config.Logger.Debug("Página escrita en caché y marcada como modificada: " + strconv.Itoa(numberPage))
			}*/

		} else { //es READ
			sizeToRead, _ := strconv.Atoi(strings.Split(instruction_string, " ")[2])
			if globals.ConfigFile.CacheEntries > 0 {
				bytesLeidosCache := readCache(pageContent, offset, sizeToRead)
				loggerString := "PID: " + strconv.Itoa(pid) + " - Acción: LEER - CACHE: " + strconv.Itoa(direccionFisica) + " - Valor: " + string(bytesLeidosCache)
				config.Logger.Info(loggerString)
			} else {
				bytesLeidos, err := ReadMemory(pid, direccionFisica, sizeToRead, globals.ConfigFile.IPMemory, globals.ConfigFile.PortMemory)
				if err != nil {
					config.Logger.Error("Error al leer memoria", "err", err)
				} else {
					config.Logger.Debug("Bytes leídos:" + string(bytesLeidos))
					loggerString := "PID: " + strconv.Itoa(pid) + " - Acción: LEER - Dirección Física: " + strconv.Itoa(direccionFisica) + " - Valor: " + string(bytesLeidos)
					config.Logger.Info(loggerString)
				}
			}

		}

	}

	// Decode de GOTO
	if len(fields) > 0 && strings.ToUpper(fields[0]) == "GOTO" {
		pcString := instruction_string[strings.Index(instruction_string, " ")+1:]
		pcInt, err := strconv.Atoi(pcString)
		if err != nil {
			config.Logger.Debug("Error al convertir PC a entero")
		}
		config.Logger.Debug("PC: " + strconv.Itoa(pcInt))
		globals.PC = pcInt
	} else {
		// Si no es GOTO, se incrementa el PC
		globals.PC += 1
	}

	// Decode DUMP_MEMORY
	if len(fields) > 0 && strings.ToUpper(fields[0]) == "DUMP_MEMORY" {
		config.Logger.Debug(fmt.Sprintf("CPU: invocando syscall DUMP_MEMORY con PID %d", pid))

		result, err := SendDumpMemoryToKernel(pid, globals.PC, globals.ConfigFile.IPKernel, globals.ConfigFile.PortKernel)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("CPU: error enviando syscall DUMP_MEMORY al Kernel: %v", err))
			// como si fallara gravemente → muere el proceso
			SendExitToKernel(pid, pc, globals.ConfigFile.IPKernel, globals.ConfigFile.PortKernel)
			return
		}

		config.Logger.Debug(fmt.Sprintf("CPU: recibió respuesta de Kernel para DUMP_MEMORY (PID %d): %s", pid, result))

		if result == "ERROR" {
			config.Logger.Info(fmt.Sprintf("CPU: proceso %d terminado por error en syscall DUMP_MEMORY", pid))
			SendExitToKernel(pid, pc, globals.ConfigFile.IPKernel, globals.ConfigFile.PortKernel)
			return
		}

		config.Logger.Info(fmt.Sprintf("CPU: syscall DUMP_MEMORY exitosa para PID %d", pid))
		return
	}

	// Decode EXIT
	if len(fields) > 0 && strings.ToUpper(fields[0]) == "EXIT" {
		err := SendExitToKernel(pid, pc, globals.ConfigFile.IPKernel, globals.ConfigFile.PortKernel)
		if err != nil {
			config.Logger.Debug("Error al enviar EXIT al kernel: " + err.Error())
		} else {
			config.Logger.Debug("Mensaje EXIT enviado al kernel con PID: " + strconv.Itoa(pid))
		}
		return
	}

	// Decode IO
	if len(fields) > 0 && strings.ToUpper(fields[0]) == "IO" {
		err := SendIOToKernel(pid, globals.PC, instruction_string, globals.ConfigFile.IPKernel, globals.ConfigFile.PortKernel)
		if err != nil {
			config.Logger.Debug("Error al enviar IO al kernel: " + err.Error())
		} else {
			config.Logger.Debug("Mensaje IO enviado al kernel con PID: " + strconv.Itoa(pid))
		}
		return

	}

	// Decode INIT_PROC

	if len(fields) > 0 && strings.ToUpper(fields[0]) == "INIT_PROC" {
		err := SendInitProcToKernel(instruction_string, globals.ConfigFile.IPKernel, globals.ConfigFile.PortKernel)
		if err != nil {
			config.Logger.Debug("Error al enviar INIT_PROC al kernel: " + err.Error())
		} else {
			config.Logger.Debug("Mensaje INIT_PROC enviado al kernel con PID: " + strconv.Itoa(pid))
		}
	}

	globals.InterruptMutex.Lock()
	//Chequear interrupciones
	if globals.InterruptExist {

		globals.InterruptMutex.Unlock()
		return
	}
	globals.InterruptMutex.Unlock()

	ProcessInstruction(pid, globals.PC)

}
