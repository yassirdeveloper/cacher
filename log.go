package main

import (
	"io"
	"log"
)

type LogType int

const (
	InfoLog LogType = iota
	WarningLog
	ErrorLog
)

type Logger interface {
	Error(string)
	Log(LogType, string)
	Warning(string)
	Info(string)
	GetLogs(*int) []string
}

type logger struct {
	*log.Logger
}

// NewLogger creates a new Logger instance with a specified prefix and flags.
func NewLogger(file io.Writer, prefix string, flags int) *logger {
	return &logger{log.New(file, prefix, flags)}
}

// Error logs a general error message.
func (l *logger) Log(logType LogType, message string) {
	switch logType {
	case InfoLog:
		l.Info(message)
	case WarningLog:
		l.Warning(message)
	case ErrorLog:
		l.Error(message)
	}
}

// Error logs a general error message.
func (l *logger) Error(message string) {
	l.Printf("%s", message)
}

// Info logs an informational message.
func (l *logger) Warning(message string) {
	l.Printf("WARNING: %s", message)
}

// Info logs an informational message.
func (l *logger) Info(message string) {
	l.Println(message)
}

func (l *logger) GetLogs(limit *int) []string {
	return []string{"test"}
}
