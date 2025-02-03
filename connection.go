package main

import (
	"bufio"
	"fmt"
	"net"
)

type Connection interface {
	Read() string
	Send(output string)
	Close()
}

type TCPConnection struct {
	net.Conn
	logger Logger
}

func NewTCPConnection(conn net.Conn, logger Logger) *TCPConnection {
	if conn == nil {
		panic("connection cannot be nil")
	}
	connection := &TCPConnection{
		Conn:   conn,
		logger: logger,
	}
	connection.logger.NewConnection(conn.RemoteAddr())
	return connection
}

func (connection *TCPConnection) Read() string {
	reader := bufio.NewReader(connection.Conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		connection.logger.Error(&UnexpectedError{message: "Error reading command", err: err})
	}
	connection.logger.Info(fmt.Sprintf("[%s] > %s", connection.RemoteAddr(), s))
	return s
}

func (connection *TCPConnection) Send(output string) {
	if connection.Conn == nil {
		connection.logger.Error(&UnexpectedError{message: "connection is not initialized", err: nil})
	}
	_, err := (*connection).Write([]byte(output))
	if err != nil {
		connection.logger.Error(&UnexpectedError{message: "error sending data", err: err})
	} else {
		connection.logger.Info(fmt.Sprintf("[%s] < %s", connection.RemoteAddr(), output))
	}
}

func (connection *TCPConnection) Close() {
	if connection.Conn == nil {
		return
	}
	err := connection.Conn.Close()
	if err != nil {
		connection.logger.Error(&UnexpectedError{message: "error closing connection", err: err})
	}
}
