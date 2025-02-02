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
)

type ServerConfig struct {
	nbrWorkers int
	port       int
}

type Server struct {
	listener    net.Listener
	config      *ServerConfig
	logger      *Logger
	connections chan Connection
	wg          sync.WaitGroup
}

var Commands = InitCommands()

func InitServer(port int, nbrWorkers int, logger *Logger) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", os.Getenv("CACHER_PORT")))
	if err != nil {
		return nil, err
	}

	return &Server{
		listener:    listener,
		config:      &ServerConfig{port: port, nbrWorkers: nbrWorkers},
		logger:      logger,
		connections: make(chan Connection),
	}, nil
}

func (server *Server) Start() {
	log.Printf("Listening on %s", server.listener.Addr())
	server.logger.ServerStarted(server.listener.Addr())

	// Create a worker pool
	for i := 0; i < server.config.nbrWorkers; i++ {
		server.wg.Add(1)
		go worker(server, server.connections)
	}

	// Accept connections and pass them to the worker pool
	go acceptConnections(server, server.connections)

	// Handle graceful shutdown
	handleShutdown(server)
}

func acceptConnections(server *Server, connections chan<- Connection) {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			server.logger.Error(&UnexpectedError{message: "Error accepting connection", err: err})
			continue
		}
		connections <- NewTCPConnection(conn, *server.logger)
	}
}

func worker(server *Server, connections <-chan Connection) {
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
	} else {
		commandInput, err := command.Parse(in[1:])
		if err != nil {
			connection.Send(err.Display())
		} else {
			output, err := command.Run(*commandInput)
			if err != nil {
				connection.Send(err.Display())
			}
			connection.Send(output.String())
		}
	}
}

func handleShutdown(server *Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		server.logger.Info("Shutting down gracefully...")
		server.listener.Close()   // Close the listener to stop accepting new connections
		close(server.connections) // Close the connections channel
		server.wg.Wait()          // Want for current connections to finish processing
	}()
}
