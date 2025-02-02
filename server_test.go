package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	// Create a logger for testing
	logTestFilePath := "server_test.log"
	logTestFile, err := os.OpenFile(logTestFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error during creating the log file: ", err)
	}
	logger := NewLogger(logTestFile, "- ", log.Ldate|log.Ltime)

	// Start the server in a goroutine
	server, err := NewServer(8080, 2, logger)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	go server.Start()

	// Wait for the server to start listening
	time.Sleep(500 * time.Millisecond) // Give the server time to start

	// Simulate a client connection
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send a command to the server
	_, err = conn.Write([]byte("TEST\n"))
	if err != nil {
		t.Fatalf("Failed to send command: %v", err)
	}

	// Simulate a shutdown signal
	time.AfterFunc(1*time.Second, func() {
		// Send a SIGINT signal to trigger shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT)
		sigChan <- syscall.SIGINT // Simulate SIGINT
	})

	// Wait for the server to shut down
	time.Sleep(2 * time.Second) // Give the server time to shut down

	// Verify that the server has shut down
	if _, err := net.Dial("tcp", "127.0.0.1:8080"); err == nil {
		t.Fatal("Server did not shut down gracefully")
	}

	// Check the logger output for graceful shutdown messages
	logOutput := strings.Join(logger.GetLogs(nil), "\n")
	if !strings.Contains(logOutput, "Shutting down gracefully") {
		t.Fatalf("Graceful shutdown message not logged. Logs: %s", logOutput)
	}
	if !strings.Contains(logOutput, "Graceful shutdown complete") {
		t.Fatalf("Graceful shutdown completion not logged. Logs: %s", logOutput)
	}
}
