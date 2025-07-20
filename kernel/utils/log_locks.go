package utils

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/config"
)

func LockLogNew(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): NEWSEM LOCKED", place, pid))
}
func UnlockLogNew(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): NEWSEM UNLOCKED", place, pid))
}

func LockLogReady(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): READYSEM LOCKED", place, pid))
}
func UnlockLogReady(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): READYSEM UNLOCKED", place, pid))
}

func LockLogSuspReady(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): SUSPREADYSEM LOCKED", place, pid))
}
func UnlockLogSuspReady(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): SUSPREADYSEM UNLOCKED", place, pid))
}

func LockLogCpu(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): CPUSEM LOCKED", place, pid))
}
func UnlockLogCpu(pid int, place string) {
	config.Logger.Debug(fmt.Sprintf("(%s) (%d): CPUSEM UNLOCKED", place, pid))
}
