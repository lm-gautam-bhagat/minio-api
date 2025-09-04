package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota + 1
	INFO
	WARN
	ERROR
	FATAL
)

var (
	currentLevel LogLevel = INFO // default level
	logger       *log.Logger
	logFile      *os.File
)

func Init(logFilePath string, level int) error {
	currentLevel = LogLevel(level)

	// Ensure log directory exists
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return err
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	logger = log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func logWithLevel(level LogLevel, levelStr string, v ...interface{}) {
	if currentLevel <= level {
		prefix := fmt.Sprintf("[%s] ", strings.ToUpper(levelStr))
		logger.Output(3, prefix+fmt.Sprintln(v...)) // 3 skips frames to point to caller
	}
}

// Debug logs a debug-level message
func Debug(v ...interface{}) {
	logWithLevel(DEBUG, "debug", v...)
}

// Info logs an info-level message
func Info(v ...interface{}) {
	logWithLevel(INFO, "info", v...)
}

// Warn logs a warning-level message
func Warn(v ...interface{}) {
	logWithLevel(WARN, "warn", v...)
}

// Error logs an error-level message
func Error(v ...interface{}) {
	logWithLevel(ERROR, "error", v...)
}

// Fatal logs a fatal-level message and exits
func Fatal(v ...interface{}) {
	logWithLevel(FATAL, "fatal", v...)
	os.Exit(1)
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
