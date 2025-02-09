package main

import (
	"bufio"
	"fmt"
	"net"
)

type Connection interface {
	Read() (string, Error)
	Send(output string) Error
	Close() Error
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
	connection.logger.Info(fmt.Sprintf("[CONNECTION_EVENT] New connection from %s", conn.RemoteAddr()))
	return connection
}

func (connection *TCPConnection) Read() (string, Error) {
	reader := bufio.NewReader(connection.Conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		err := &UnexpectedError{message: "Error reading command", err: err}
		connection.logger.Error(fmt.Sprintf("[CONNECTION_EVENT] %s", err.Error()))
		return "", err
	}
	connection.logger.Info(fmt.Sprintf("[CONNECTION_EVENT] [%s] > %s", connection.RemoteAddr(), s))
	return s, nil
}

func (connection *TCPConnection) Send(output string) Error {
	if connection.Conn == nil {
		err := &UnexpectedError{message: "connection is not initialized", err: nil}
		connection.logger.Error(fmt.Sprintf("[CONNECTION_EVENT] %s", err.Error()))
		return err
	}
	_, err := (*connection).Write([]byte(output))
	if err != nil {
		err := &UnexpectedError{message: "error sending data", err: err}
		connection.logger.Error(fmt.Sprintf("[CONNECTION_EVENT] %s", err.Error()))
		return err
	}
	connection.logger.Info(fmt.Sprintf("[CONNECTION_EVENT] [%s] < %s", connection.RemoteAddr(), output))
	return nil
}

func (connection *TCPConnection) Close() Error {
	if connection.Conn == nil {
		err := &UnexpectedError{message: "tried closing inexistant connection", err: nil}
		connection.logger.Error(fmt.Sprintf("[CONNECTION_EVENT] %s", err.Error()))
		return err
	}
	err := connection.Conn.Close()
	if err != nil {
		err := &UnexpectedError{message: "error closing connection", err: err}
		connection.logger.Error(fmt.Sprintf("[CONNECTION_EVENT] %s", err.Error()))
		return err
	}
	return nil
}
