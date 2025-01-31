package main

import (
	"net"
	"reflect"
)

type CommandArgument struct {
	position  int
	valueType reflect.Type
}

// if the option does not need a value it shall have a valueType of nil
type CommandOption struct {
	letter    rune
	valueType reflect.Type
}

type Command struct {
	arguments map[string]*CommandArgument
	options   map[string]*CommandOption
	write     bool
	handler   func(args []string) ([]byte, error)
}

type Request struct {
	*Command
	connection net.Conn
}

var GetCommand *Command = &Command{
	arguments: map[string]*CommandArgument{
		"key": {position: 0, valueType: reflect.TypeOf("s")},
	},
	options: map[string]*CommandOption{},
	handler: getHandler,
}

func getHandler(args []string) ([]byte, error) {
	return []byte("success"), nil
}
