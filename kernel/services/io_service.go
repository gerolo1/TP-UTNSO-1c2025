package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/clients"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/config"
	"github.com/sisoputnfrba/tp-golang/kernel/models"
)

var IoDevices = map[string]*models.IoDevice{
	"DUMP": {
		Queue:     []*models.IoQueueEntry{},
		Instances: []*models.IoInstance{},
	},
}
var IoDevicesSem sync.Mutex

// Registrar un dispositivo
func RegisterIoDevice(name string, ip string, port int) {
	IoDevicesSem.Lock()

	device, exists := IoDevices[name]
	if !exists {
		device = &models.IoDevice{
			Queue:     []*models.IoQueueEntry{},
			Instances: []*models.IoInstance{},
		}
		IoDevices[name] = device
	}

	instance := &models.IoInstance{
		Url:     fmt.Sprintf("http://%s:%d", ip, port),
		Process: nil, // free
	}

	device.Instances = append(device.Instances, instance)

	IoDevicesSem.Unlock()

	config.Logger.Info(fmt.Sprintf("Registered IO Device: %s at %s:%d", name, ip, port))
}

func StartTimerSuspend(pid int, device *models.IoDevice) {
	time.Sleep(time.Duration(config.Config.SuspensionTime) * time.Millisecond)
	NewSem.Lock()
	SuspReadySem.Lock()
	ReadySem.Lock()
	CpuSem.Lock()
	IoDevicesSem.Lock()
	ioEntry := FindIoEntryInIoDeviceByPid(pid, device)
	if ioEntry == nil {
		instance := FindIoInstanceByPid(pid, device)
		if instance != nil {
			ioEntry = instance.Process
		}
	}
	if ioEntry != nil {

		UpdatePcb(ioEntry.Pcb, "SUSPBLOCKED")

		config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado BLOCKED al estado SUSPBLOCKED", pid))

		clients.SuspendProcess(pid)

		config.Logger.Debug(fmt.Sprintf("FINISHED SUSPENDING PID %d", pid))

		enoughMemory := true
		for len(SuspReadyList) > 0 && enoughMemory {
			enoughMemory = CheckSuspReady()
		}
		for len(NewList) > 0 && enoughMemory {
			enoughMemory = CheckNew()
		}
		ReadyToExec()
	}
	NewSem.Unlock()
	SuspReadySem.Unlock()
	ReadySem.Unlock()
	CpuSem.Unlock()
	IoDevicesSem.Unlock()
	config.Logger.Debug("DESBLOQUEADOS LOS SEMAFOROS DE STARTTIMERSUSPEND")
}

func AddToBlocked(device *models.IoDevice, pc int, duration int, pcb *models.Pcb, deviceName string) {
	pcb.Pc = pc

	entry := &models.IoQueueEntry{
		Pcb:      pcb,
		Duration: duration,
	}
	device.Queue = append(device.Queue, entry)

	UpdatePcb(pcb, "BLOCKED")

	config.Logger.Info(fmt.Sprintf("## (%d) - Bloqueado por IO: %s", pcb.Pid, deviceName))

	go StartTimerSuspend(pcb.Pid, device)
}

// Manejo de Syscall
func HandleIoSyscall(pid int, pc int, deviceName string, duration int) {
	NewSem.Lock()
	utils.LockLogNew(pid, "HANDLEIOSYSCALL")
	SuspReadySem.Lock()
	utils.LockLogSuspReady(pid, "HANDLEIOSYSCALL")
	ReadySem.Lock()
	utils.LockLogReady(pid, "HANDLEIOSYSCALL")
	CpuSem.Lock()
	utils.LockLogCpu(pid, "HANDLEIOSYSCALL")
	IoDevicesSem.Lock()
	config.Logger.Debug(fmt.Sprintf(" (%d) IODEVICES LOCK", pid))

	config.Logger.Debug(fmt.Sprintf("INIT HandleIoSyscall - PID: %d, Device: %s, Duration: %d", pid, deviceName, duration))
	device, exists := IoDevices[deviceName]

	if !exists {
		config.Logger.Debug(fmt.Sprintf("Device %s does not exist. Sending process %d to EXIT.", deviceName, pid))
		ExitProcessByPid(pid)
		NewSem.Unlock()
		utils.UnlockLogNew(pid, "HANDLEIOSYSCALL")
		SuspReadySem.Unlock()
		utils.UnlockLogSuspReady(pid, "HANDLEIOSYSCALL")
		ReadySem.Unlock()
		utils.UnlockLogReady(pid, "HANDLEIOSYSCALL")
		CpuSem.Unlock()
		utils.UnlockLogCpu(pid, "HANDLEIOSYSCALL")
		IoDevicesSem.Unlock()
		return
	}

	NewSem.Unlock()
	utils.UnlockLogNew(pid, "HANDLEIOSYSCALL")
	SuspReadySem.Unlock()
	utils.UnlockLogSuspReady(pid, "HANDLEIOSYSCALL")

	// Buscar el CPU
	cpu := FindCpuByPid(pid)
	if cpu == nil {
		config.Logger.Debug("LA CPU ME DA NULL")
	}

	AddToBlocked(device, pc, duration, cpu.Pcb, deviceName)
	config.Logger.Debug(fmt.Sprintf("Process PID %d enqueued in IO device %s with duration %d", pid, deviceName, duration))

	cpu.Busy = false

	if len(ReadyList) > 0 {
		ReadyToExec()
	}

	instance := FindFreeIoInstance(device)
	// Si está libre, arrancar el primero
	if instance != nil {
		assignNextProcessToIo(device, instance, deviceName)
	}

	ReadySem.Unlock()
	utils.UnlockLogReady(pid, "HANDLEIOSYSCALL")
	CpuSem.Unlock()
	utils.UnlockLogCpu(pid, "HANDLEIOSYSCALL")
	IoDevicesSem.Unlock()
	config.Logger.Debug(fmt.Sprintf(" (%d) IODEVICES UNLOCK", pid))

	config.Logger.Debug(fmt.Sprintf("END HandleIoSyscall - PID: %d", pid))
}

// Asigo al prox en la cola

func assignNextProcessToIo(device *models.IoDevice, instance *models.IoInstance, deviceName string) {
	if len(device.Queue) == 0 {
		return
	}

	entry := device.Queue[0]
	device.Queue = slices.Delete(device.Queue, 0, 1)
	instance.Process = entry

	config.Logger.Info(fmt.Sprintf("## (%d) realiza IO %s, Instancia: %s, Duracion: %d", entry.Pcb.Pid, deviceName, instance.Url, entry.Duration))

	sendIoRequestToModule(instance, deviceName)
}

// Mando a IO

func sendIoRequestToModule(instance *models.IoInstance, deviceName string) {
	url := fmt.Sprintf("%s/io/request", instance.Url)

	payload := map[string]interface{}{
		"pid":      instance.Process.Pcb.Pid,
		"duration": instance.Process.Duration,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		config.Logger.Error(fmt.Sprintf("Failed to serialize IO request for PID %d: %v", instance.Process.Pcb.Pid, err))
		markInstanceFree(instance)
		return
	}

	config.Logger.Debug(fmt.Sprintf("Sending IO request to %s - Payload: %s", url, string(jsonData)))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		config.Logger.Error(fmt.Sprintf("Failed to send IO request to %s: %v", url, err))
		markInstanceFree(instance)
		return
	}
	defer resp.Body.Close()

	config.Logger.Debug(fmt.Sprintf("IO request sent to module %s for PID %d", deviceName, instance.Process.Pcb.Pid))
}

func BlockedToReady(pcb *models.Pcb) {
	UpdatePcb(pcb, "READY")
	AddToReady(pcb)
	config.Logger.Info(fmt.Sprintf("## (%d) finalizo IO y pasa de BLOCKED a READY", pcb.Pid))
	if len(ReadyList) > 0 {
		ReadyToExec()
	}
}

func SuspBlockedToSuspReady(pcb *models.Pcb) {
	UpdatePcb(pcb, "SUSPREADY")
	SuspReadyList = append(SuspReadyList, pcb)

	config.Logger.Info(fmt.Sprintf("## (%d) finalizo IO y pasa de SUSPBLOCKED a SUSPREADY", pcb.Pid))

	var tryToExecute bool = true
	switch config.Config.ReadyIngressAlgorithm {
	case "FIFO":
		if len(SuspReadyList) == 1 {
			tryToExecute = true
		}
	case "PMCP":
		slices.SortFunc(SuspReadyList, PMCPAlgorithm)
		if SuspReadyList[0].Pid == pcb.Pid {
			tryToExecute = true
		}
	}
	if tryToExecute {
		response := clients.ResumeProcess(pcb)

		if response.StatusCode == http.StatusOK {
			RemoveFromCurrentStateList(pcb)

			UpdatePcb(pcb, "READY")

			AddToReady(pcb)
			config.Logger.Info(fmt.Sprintf("## (%d) Pasa del estado SUSP. READY al estado READY", pcb.Pid))

			ReadyToExec()
		} else {
			config.Logger.Debug(fmt.Sprintf("(%d) NOT INITIALIZED: NOT ENOUGH MEMORY", pcb.Pid))
		}
	}
}

// ---------- HANDLE IO COMPLETE ----------
func HandleIoComplete(deviceName string, pid int) {
	SuspReadySem.Lock()
	utils.LockLogSuspReady(pid, "IOCOMPLETE")
	ReadySem.Lock()
	utils.LockLogReady(pid, "IOCOMPLETE")
	CpuSem.Lock()
	utils.LockLogCpu(pid, "IOCOMPLETE")
	IoDevicesSem.Lock()
	config.Logger.Debug(fmt.Sprintf(" (%d) IODEVICES LOCK", pid))

	config.Logger.Debug(fmt.Sprintf("INIT HandleIoComplete - Device: %s, PID: %d", deviceName, pid))

	// Buscar device
	device, exists := IoDevices[deviceName]
	instance := FindIoInstanceByPid(pid, device)

	if !exists {
		config.Logger.Error(fmt.Sprintf("Device %s not found on completion", deviceName))
		return
	}

	pcb := instance.Process.Pcb

	// Marcar libre
	markInstanceFree(instance)
	// Verificar siguiente en cola

	assignNextProcessToIo(device, instance, deviceName)

	IoDevicesSem.Unlock()
	config.Logger.Debug(fmt.Sprintf(" (%d) IODEVICES UNLOCK", pid))

	if pcb.CurrentState == "BLOCKED" {
		SuspReadySem.Unlock()
		utils.UnlockLogSuspReady(pid, "IOCOMPLETE")

		BlockedToReady(pcb)
	} else {
		SuspBlockedToSuspReady(pcb)

		SuspReadySem.Unlock()
		utils.UnlockLogSuspReady(pid, "IOCOMPLETE")
	}
	ReadySem.Unlock()
	utils.UnlockLogReady(pid, "IOCOMPLETE")
	CpuSem.Unlock()
	utils.UnlockLogCpu(pid, "IOCOMPLETE")
}

func RemoveQueueIfNoInstances(device *models.IoDevice, memoryFreed bool) bool {
	if len(device.Instances) == 0 {
		for len(device.Queue) > 0 {
			pcb := device.Queue[0].Pcb
			device.Queue = slices.Delete(device.Queue, 0, 1)

			if pcb.CurrentState == "BLOCKED" {
				memoryFreed = true
			}

			clients.ExitMemory(pcb.Pid)

			state := pcb.CurrentState
			UpdatePcb(pcb, "EXIT")

			config.Logger.Info(fmt.Sprintf("## (%d) - Finaliza el proceso desde estado %s", pcb.Pid, state))

			newState := pcb.StateMap["NEW"]
			ready := pcb.StateMap["READY"]
			exec := pcb.StateMap["EXEC"]
			blocked := pcb.StateMap["BLOCKED"]
			suspBlocked := pcb.StateMap["SUSPBLOCKED"]
			suspReady := pcb.StateMap["SUSPREADY"]
			exit := pcb.StateMap["EXIT"]
			config.Logger.Info(fmt.Sprintf("## (%d) - Métricas de estado: NEW (%d) (%d), READY (%d) (%d), EXEC (%d) (%d), BLOCKED (%d) (%d), SUSP. BLOCKED (%d) (%d), SUSP. READY (%d) (%d), EXIT (%d) (%d)",
				pcb.Pid, newState.Count, newState.Time, ready.Count, ready.Time, exec.Count, exec.Time, blocked.Count, blocked.Time, suspBlocked.Count, suspBlocked.Time, suspReady.Count, suspReady.Time, exit.Count, exit.Time))
		}
	}
	return memoryFreed
}

// si hay procesos y no hay instancias pasamos todos a exit
// si no hay procesos y tampoco instancias, no hay que hacer nada
func HandleIoShutdown(deviceName string, ip string, port int) {
	NewSem.Lock()
	utils.LockLogNew(101, "HANDLEIOSHUTDOWN")
	SuspReadySem.Lock()
	utils.LockLogSuspReady(101, "HANDLEIOSHUTDOWN")
	ReadySem.Lock()
	utils.LockLogReady(101, "HANDLEIOSHUTDOWN")
	CpuSem.Lock()
	utils.LockLogCpu(101, "HANDLEIOSHUTDOWN")
	IoDevicesSem.Lock()
	config.Logger.Debug("IODEVICESLOCK IOSHUTDOWN")

	device, _ := IoDevices[deviceName]
	instance := FindIoInstanceByUrl(device, ip, port)
	var memoryFreed bool = false

	if instance.Process != nil {
		if instance.Process.Pcb.CurrentState == "BLOCKED" {
			memoryFreed = true
		}

		clients.ExitMemory(instance.Process.Pcb.Pid)

		state := instance.Process.Pcb.CurrentState
		UpdatePcb(instance.Process.Pcb, "EXIT")

		config.Logger.Info(fmt.Sprintf("## (%d) - Finaliza el proceso desde estado %s", instance.Process.Pcb.Pid, state))

		newState := instance.Process.Pcb.StateMap["NEW"]
		ready := instance.Process.Pcb.StateMap["READY"]
		exec := instance.Process.Pcb.StateMap["EXEC"]
		blocked := instance.Process.Pcb.StateMap["BLOCKED"]
		suspBlocked := instance.Process.Pcb.StateMap["SUSPBLOCKED"]
		suspReady := instance.Process.Pcb.StateMap["SUSPREADY"]
		exit := instance.Process.Pcb.StateMap["EXIT"]
		config.Logger.Info(fmt.Sprintf("## (%d) - Métricas de estado: NEW (%d) (%d), READY (%d) (%d), EXEC (%d) (%d), BLOCKED (%d) (%d), SUSP. BLOCKED (%d) (%d), SUSP. READY (%d) (%d), EXIT (%d) (%d)",
			instance.Process.Pcb.Pid, newState.Count, newState.Time, ready.Count, ready.Time, exec.Count, exec.Time, blocked.Count, blocked.Time, suspBlocked.Count, suspBlocked.Time, suspReady.Count, suspReady.Time, exit.Count, exit.Time))
	}
	instanceIndex := GetIoInstanceIndex(device, ip, port)
	device.Instances = slices.Delete(device.Instances, instanceIndex, instanceIndex+1)

	memoryFreed = RemoveQueueIfNoInstances(device, memoryFreed)

	if len(device.Instances) == 0 {
		delete(IoDevices, deviceName)
	}

	if memoryFreed {
		var enoughMemory bool = true
		for len(SuspReadyList) > 0 && enoughMemory {
			enoughMemory = CheckSuspReady()
		}
		for len(NewList) > 0 && enoughMemory {
			enoughMemory = CheckNew()
		}
		ReadyToExec()
	}
	NewSem.Unlock()
	SuspReadySem.Unlock()
	ReadySem.Unlock()
	CpuSem.Unlock()
	IoDevicesSem.Unlock()
}

// ---------- UTILS ----------

func GetIoInstanceIndex(device *models.IoDevice, ip string, port int) int {
	url := fmt.Sprintf("http://%s:%d", ip, port)
	for index, instance := range device.Instances {
		if instance.Url == url {
			return index
		}
	}
	return 0
}

func FindIoInstanceByUrl(device *models.IoDevice, ip string, port int) *models.IoInstance {
	url := fmt.Sprintf("http://%s:%d", ip, port)
	for _, instance := range device.Instances {
		if instance.Url == url {
			return instance
		}
	}
	return nil
}

func FindIoInstanceByPid(pid int, device *models.IoDevice) *models.IoInstance {
	for _, instance := range device.Instances {
		if instance.Process != nil && instance.Process.Pcb.Pid == pid {
			return instance
		}
	}
	return nil
}

func FindIoEntryInIoDeviceByPid(pid int, device *models.IoDevice) *models.IoQueueEntry {
	for _, ioEntry := range device.Queue {
		if ioEntry.Pcb.Pid == pid {
			return ioEntry
		}
	}
	return nil
}

func FindFreeIoInstance(device *models.IoDevice) *models.IoInstance {
	for _, ioInstance := range device.Instances {
		if ioInstance.Process == nil {
			return ioInstance
		}
	}
	return nil
}

func markInstanceFree(instance *models.IoInstance) {
	instance.Process = nil
}

func ExitProcessByPid(pid int) {
	config.Logger.Debug(fmt.Sprintf("INIT ExitProcessByPid - PID: %d", pid))

	clients.ExitMemory(pid)

	cpu := FindCpuByPid(pid)

	ExecToExit(cpu)

	if len(SuspReadyList) == 0 && len(NewList) == 0 && len(ReadyList) > 0 {
		ReadyToExec()
	} else {
		enoughMemory := true
		for len(SuspReadyList) > 0 && enoughMemory {
			enoughMemory = CheckSuspReady()
		}
		for len(NewList) > 0 && enoughMemory {
			enoughMemory = CheckNew()
		}
		ReadyToExec()
	}

	config.Logger.Info(fmt.Sprintf("## (%d) - Proceso finalizado por error en syscall IO", pid))
}
