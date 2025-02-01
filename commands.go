package main

import (
	"reflect"
)

type CommandArgument struct {
	position  int
	valueType reflect.Type
}

// if the option does not need a value it shall have a valueType of nil
type CommandOption struct {
	letter    rune
	name      string
	valueType reflect.Type
}

type Command struct {
	arguments map[string]*CommandArgument
	options   map[string]*CommandOption
	write     bool
	handler   func(args []string) ([]byte, error)
}

func (command *Command) Run(in []string) (string, Error) {
	return "OK", nil
}

var GetCommand *Command = &Command{
	arguments: map[string]*CommandArgument{
		"key": {position: 0, valueType: reflect.TypeOf(string(""))},
	},
	options: map[string]*CommandOption{
		"expiration": {letter: 'e', name: "expires", valueType: reflect.TypeOf(int(0))},
	},
	handler: func(in []string) ([]byte, error) {
		return []byte("success"), nil
	},
}
