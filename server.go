package main

import (
	"bufio"
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
	commands    map[string]*Command
	logger      *Logger
	connections chan net.Conn
	wg          sync.WaitGroup
}

func InitServer(port int, nbrWorkers int, logger *Logger) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", os.Getenv("CACHER_PORT")))
	if err != nil {
		return nil, err
	}

	commands := map[string]*Command{
		"GET": GetCommand,
	}

	return &Server{
		listener:    listener,
		config:      &ServerConfig{port: port, nbrWorkers: nbrWorkers},
		commands:    commands,
		logger:      logger,
		connections: make(chan net.Conn),
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

	// Wait for workers to finish
	server.wg.Wait()
}

func acceptConnections(server *Server, connections chan<- net.Conn) {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			server.logger.Error(&UnexpectedError{message: "Error accepting connection", err: err})
			continue
		}
		server.logger.NewConnection(conn.RemoteAddr())
		connections <- conn
	}
}

func worker(server *Server, connections <-chan net.Conn) {
	for conn := range connections {
		server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(connection net.Conn) {
	defer connection.Close()

	reader := bufio.NewReader(connection)
	commandString, err := reader.ReadString('\n')
	if err != nil {
		server.logger.Error(&UnexpectedError{message: "Error reading command", err: err})
		return
	}
	server.logger.Info(fmt.Sprintf("[%s] > %s", connection.RemoteAddr(), commandString))

	l := strings.Split(commandString, " ")
	commandName := strings.ToUpper(l[0])
	command := server.commands[commandName]
	if command == nil {
		_, err = connection.Write([]byte("invalid command"))
		if err != nil {
			server.logger.Error(&UnexpectedError{message: "Error writing to connection", err: err})
			return
		}
		server.logger.Error(&InvalidCommandError{command: commandName})
	} else {
		request := &Request{command, connection}
		if command.write {
			handleWriteRequest(request)
		} else {
			handleReadRequest(request)
		}
	}
}

func handleReadRequest(request *Request) {

}

func handleWriteRequest(request *Request) {

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
