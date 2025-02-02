package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ServerConfig struct {
	nbrWorkers int
	port       int
}

type Server struct {
	listener    net.Listener
	config      *ServerConfig
	logger      Logger
	connections chan Connection
	shutdown    chan os.Signal
	wg          sync.WaitGroup
}

var Commands = InitCommands()

func NewServer(port int, nbrWorkers int, logger Logger) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}

	return &Server{
		listener:    listener,
		config:      &ServerConfig{port: port, nbrWorkers: nbrWorkers},
		logger:      logger,
		shutdown:    nil,
		connections: make(chan Connection),
	}, nil
}

func (server *Server) Start() {
	log.Printf("Listening on %s", server.listener.Addr().String())
	server.logger.Info(fmt.Sprintf("Server started on %s", server.listener.Addr()))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	server.shutdown = sigChan

	// Handle graceful shutdown
	go handleShutdown(server)

	// Create a worker pool
	for i := 0; i < server.config.nbrWorkers; i++ {
		server.wg.Add(1)
		go worker(server, server.connections)
	}

	// Accept connections and pass them to the worker pool
	go acceptConnections(server, server.connections)

	server.wg.Wait()
}

func acceptConnections(server *Server, connections chan<- Connection) {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			// Check if the listener is closed (graceful shutdown)
			select {
			case <-server.shutdown: // If shutdown is signaled, exit the loop
				server.logger.Info("Connection acceptance stopped")
				return
			default:
				if !isClosedConnectionError(err) {
					server.logger.Error(&UnexpectedError{message: "Error accepting connection", err: err})
				}
			}
		} else {
			connections <- NewTCPConnection(conn, server.logger)
		}
	}
}

func worker(server *Server, connections <-chan Connection) {
	defer server.wg.Done()

	for conn := range connections {
		server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(connection Connection) {
	defer connection.Close()

	commandString := connection.Read()
	in := strings.Split(strings.TrimSpace(commandString), " ")
	commandName := strings.ToUpper(in[0])
	command := Commands[commandName]
	if command == nil {
		invalidCommandError := &InvalidCommandError{command: commandName}
		connection.Send(invalidCommandError.Display())
		return
	}
	commandInput, err := command.Parse(in[1:])
	if err != nil {
		connection.Send(err.Display())
		return
	}
	output, err := command.Run(*commandInput)
	if err != nil {
		connection.Send(err.Display())
		return
	}
	connection.Send(output.String())
}

func handleShutdown(server *Server) {
	<-server.shutdown // Wait for a shutdown signal

	server.logger.Info("Shutting down gracefully...")
	server.listener.Close()   // Stop accepting new connections
	close(server.connections) // Close the connections channel

	// Wait for workers to finish with a timeout
	done := make(chan struct{})
	go func() {
		server.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		server.logger.Info("Graceful shutdown complete")
	case <-time.After(5 * time.Second):
		server.logger.Warning("Forcing shutdown after timeout")
	}
}

func isClosedConnectionError(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
