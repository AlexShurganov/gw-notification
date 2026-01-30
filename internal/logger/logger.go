package logger

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
)

func InitLogger(service string) (*slog.Logger, error) {
	log.Println("Initializing logger...")
	file, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), nil))

	logger.Debug(fmt.Sprintf("Service %s===== Debug message", service))
	logger.Info(fmt.Sprintf("Service %s===== Info message", service))
	logger.Warn(fmt.Sprintf("Service %s===== Warning message", service))
	logger.Error(fmt.Sprintf("Service %s===== Error message", service))

	return logger, nil
}
