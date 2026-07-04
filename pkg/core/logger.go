package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelWarning
	LevelInfo
	LevelDebug
)

var (
	GlobalLogLevel LogLevel
	logger         *log.Logger
)

func init() {
	// 默认禁用日志输出，除非显式初始化
	logger = log.New(io.Discard, "", 0)
	GlobalLogLevel = LevelError
}

func parseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warning":
		return LevelWarning
	case "error":
		return LevelError
	default:
		return LevelError
	}
}

func InitLogger(level string, logFile string) {
	GlobalLogLevel = parseLevel(level)

	var output io.Writer = os.Stdout
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			output = file
		} else {
			fmt.Printf("Failed to open log file %s: %v\n", logFile, err)
		}
	}

	// 在非 debug 模式下，如果不指定 loglevel 且不指定 logfile，不干扰默认输出流
	if level == "" && logFile == "" {
		output = io.Discard
	}

	logger = log.New(output, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func LogDebug(format string, v ...any) {
	if GlobalLogLevel >= LevelDebug {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

func LogInfo(format string, v ...any) {
	if GlobalLogLevel >= LevelInfo {
		logger.Printf("[INFO] "+format, v...)
	}
}

func LogWarning(format string, v ...any) {
	if GlobalLogLevel >= LevelWarning {
		logger.Printf("[WARNING] "+format, v...)
	}
}

func LogError(format string, v ...any) {
	if GlobalLogLevel >= LevelError {
		logger.Printf("[ERROR] "+format, v...)
	}
}
