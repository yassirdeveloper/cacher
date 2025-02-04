package main

import (
	"bufio"
	"fmt"
	"net"
)

type Connection interface {
	Log(LogType, string)
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

func (connection *TCPConnection) Log(logType LogType, message string) {
	connection.logger.Log(logType, message)
}

func (connection *TCPConnection) Read() string {
	reader := bufio.NewReader(connection.Conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		err := &UnexpectedError{message: "Error reading command", err: err}
		connection.Log(ErrorLog, err.Error())
	}
	connection.logger.Info(fmt.Sprintf("[%s] > %s", connection.RemoteAddr(), s))
	return s
}

func (connection *TCPConnection) Send(output string) {
	if connection.Conn == nil {
		err := &UnexpectedError{message: "connection is not initialized", err: nil}
		connection.Log(ErrorLog, err.Error())
	}
	_, err := (*connection).Write([]byte(output))
	if err != nil {
		err := &UnexpectedError{message: "error sending data", err: err}
		connection.Log(ErrorLog, err.Error())
	} else {
		connection.Log(InfoLog, fmt.Sprintf("[%s] < %s", connection.RemoteAddr(), output))
	}
}

func (connection *TCPConnection) Close() {
	if connection.Conn == nil {
		return
	}
	err := connection.Conn.Close()
	if err != nil {
		err := &UnexpectedError{message: "error closing connection", err: err}
		connection.Log(ErrorLog, err.Error())
	}
}
