package main

import (
	"fmt"
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

type Server interface {
	Start(time.Duration)
	acceptConnection() (Connection, Error)
	handleConnection(Connection)
	CloseConnections()
	Log(LogType, string)
	wait()
	done()
	ShutDown(time.Duration)
	ShutDownChan() chan os.Signal
}

type server[K comparable, V any] struct {
	listener       net.Listener
	config         *ServerConfig
	logger         Logger
	connections    chan Connection
	commandManager CommandManager
	cacheManager   CacheManager[K, V]
	executor       Executor[K, V]
	shutdown       chan os.Signal
	wg             sync.WaitGroup
}

func NewServer[K comparable, V any](port int, nbrWorkers int, logger Logger, commandManager CommandManager, cacheManager CacheManager[K, V]) (Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}

	return &server[K, V]{
		listener:       listener,
		config:         &ServerConfig{port: port, nbrWorkers: nbrWorkers},
		logger:         logger,
		shutdown:       nil,
		connections:    make(chan Connection),
		commandManager: commandManager,
		cacheManager:   cacheManager,
		executor:       NewExecutor(cacheManager),
	}, nil
}

func (server *server[K, V]) acceptConnection() (Connection, Error) {
	conn, err := server.listener.Accept()
	if err != nil {
		// Check if the listener is closed (graceful shutdown)
		select {
		case <-server.shutdown: // If shutdown is signaled, exit the loop
			server.Log(InfoLog, "Connection acceptance stopped")
			return nil, nil
		default:
			if !isClosedConnectionError(err) {
				return nil, &UnexpectedError{message: "Error accepting connection", err: err}
			}
			return nil, nil
		}
	}
	return NewTCPConnection(conn, server.logger), nil
}

func (server *server[K, V]) handleConnection(connection Connection) {
	defer connection.Close()

	commandString, err := connection.Read()
	if err != nil {
		connection.Send(err.Display())
		return
	}
	in := strings.Split(strings.TrimSpace(commandString), " ")
	commandName := strings.ToUpper(in[0])
	command, err := server.commandManager.Get(commandName)
	if err != nil {
		connection.Send(err.Display())
		return
	}
	commandInput, err := command.Parse(in[1:])
	if err != nil {
		connection.Send(err.Display())
		return
	}
	executableCommand, ok := command.(ExecutableCommand[K, V])
	if ok {
		result, err := server.executor.Execute(executableCommand, commandInput)
		if err != nil {
			connection.Send(err.Display())
		} else {
			connection.Send(result.String())
		}
	} else {
		err := &CommandNotExecutableError{command: commandName}
		connection.Send(err.Display())
	}
}

func (server *server[K, V]) CloseConnections() {
	server.listener.Close()   // Stop accepting new connections
	close(server.connections) // Close the connections channel
}

func (server *server[K, V]) wait() {
	server.wg.Wait()
}

func (server *server[K, V]) done() {
	server.wg.Done()
}

func (server *server[K, V]) Log(logType LogType, message string) {
	server.logger.Log(logType, message)
}

func (server *server[K, V]) Start(timeout time.Duration) {
	server.Log(InfoLog, fmt.Sprintf("Server started on %s", server.listener.Addr()))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	server.shutdown = sigChan

	// Handle graceful shutdown
	go handleShutdown(server, timeout)

	// Create a worker pool
	for i := 0; i < server.config.nbrWorkers; i++ {
		server.wg.Add(1)
		go worker(server, server.connections)
	}

	// Accept connections and pass them to the worker pool
	go acceptConnections(server, server.connections)

	server.wg.Wait()
}

func (server *server[K, V]) ShutDownChan() chan os.Signal {
	return server.shutdown
}

func (server *server[K, V]) ShutDown(timeout time.Duration) {
	server.Log(InfoLog, "Shutting down gracefully...")
	server.CloseConnections()
	// Wait for workers to finish with a timeout
	done := make(chan struct{})
	go func() {
		server.wait()
		close(done)
	}()

	select {
	case <-done:
		server.Log(InfoLog, "Graceful shutdown complete")
	case <-time.After(timeout):
		server.Log(WarningLog, "Forcing shutdown after timeout")
	}
}

func acceptConnections(server Server, connections chan<- Connection) {
	for {
		connection, err := server.acceptConnection()
		if err != nil {
			server.Log(ErrorLog, err.Error())
		} else {
			connections <- connection
		}
	}
}

func worker(server Server, connections <-chan Connection) {
	defer server.done()

	for conn := range connections {
		server.handleConnection(conn)
	}
}

func handleShutdown(server Server, timeout time.Duration) {
	<-server.ShutDownChan() // Wait for a shutdown signal
	server.ShutDown(timeout)
}

func isClosedConnectionError(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
