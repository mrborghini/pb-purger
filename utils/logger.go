package utils

import (
	"fmt"
	"time"
)

type LogLevel string
type Color string

const (
	Info    LogLevel = "\033[36m[INFO"
	Warning LogLevel = "\033[33m[WARNING"
	Error   LogLevel = "\033[31m[ERROR"
	Debug   LogLevel = "\033[32m[DEBUG"
)

type Logger struct {
	typeName string
}

// Initialize a new logger
func NewLogger(typeName string) *Logger {
	return &Logger{typeName: typeName}
}

// Log a message
func (l *Logger) log(loglevel LogLevel, message string) {
	var shownMessage string

	shownMessage = fmt.Sprintf("%s - %s - %s]: %s\033[m", loglevel, l.typeName, getCurrentTime(), message)

	fmt.Println(shownMessage)
}

// Log an info message
func (l *Logger) Info(message string) {
	l.log(Info, message)
}

// Log a warning message
func (l *Logger) Warning(message string) {
	l.log(Warning, message)
}

// Log an error message
func (l *Logger) Error(message string) {
	l.log(Error, message)
}

// Log a debug message
func (l *Logger) Debug(message string) {
	l.log(Debug, message)
}

func getCurrentTime() string {
	dt := time.Now()
	return dt.Local().Format("2006-01-02 15:04:05")
}
