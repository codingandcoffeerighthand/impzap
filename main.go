package main

import (
	"doitsolutions.vn/pkg/configs"
	"doitsolutions.vn/pkg/logger"
)

func main() {
	logCfg := configs.LogConfig{
		LogLevel:    "debug",
		Dir:         "logs",
		MaxBackups:  5,
		MaxSize:     10,
		MaxAge:      30,
		Compress:    true,
		ShowConsole: true,
	}
	logger := logger.New(logCfg)
	defer logger.Sync()

	// Example usage
	logger.Debug("This debug message will only appear in console")
	logger.Info("This is an info log")
	logger.Warn("This warning will appear in console and error log")
	logger.Error("This is an error log")
}
