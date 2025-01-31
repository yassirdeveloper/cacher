package main

import (
	"io"
	"log"
	"net"
)

type Logger struct {
	*log.Logger
}

// NewLogger creates a new Logger instance with a specified prefix and flags.
func NewLogger(writer io.Writer, prefix string, flags int) *Logger {
	return &Logger{log.New(writer, prefix, flags)}
}

// NewConnection logs a message when the server starts.
func (logger *Logger) ServerStarted(addr net.Addr) {
	logger.Printf("Server started on %s", addr)
}

// NewConnection logs a message for a new connection.
func (logger *Logger) NewConnection(addr net.Addr) {
	logger.Printf("New connection from %s", addr)
}

// Error logs a general error message.
func (l *Logger) Error(err error) {
	l.Printf("%s", err)
}

// Info logs an informational message.
func (l *Logger) Info(message string) {
	l.Println(message)
}
