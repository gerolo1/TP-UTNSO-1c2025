package utils

import (
	"log/slog"
	"strconv"
	"time"
)

func RequestDelay(logger *slog.Logger, delay int) { //ToDo: Revisar donde se usa esto
	logger.Debug("INIT: DelayUtils - RequestDelay / Delay: " + strconv.Itoa(delay) + "ms")

	time.Sleep(time.Duration(delay) * time.Millisecond)

	logger.Debug("END: DelayUtils - RequestDelay / Delay: " + strconv.Itoa(delay) + "ms")
}
