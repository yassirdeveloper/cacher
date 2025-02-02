package main

import (
	"io"
	"log"
	"net"
)

type Logger interface {
	Error(err Error)
	Warning(message string)
	Info(message string)
	NewConnection(addr net.Addr)
	GetLogs(limit *int) []string
}

type logger struct {
	*log.Logger
}

// NewLogger creates a new Logger instance with a specified prefix and flags.
func NewLogger(file io.Writer, prefix string, flags int) *logger {
	return &logger{log.New(file, prefix, flags)}
}

// NewConnection logs a message for a new connection.
func (logger *logger) NewConnection(addr net.Addr) {
	logger.Printf("New connection from %s", addr)
}

// Error logs a general error message.
func (l *logger) Error(err Error) {
	l.Printf("%s", err)
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
