package main

import (
	"slices"
)

const (
	TypeInt    = "int"
	TypeFloat  = "float"
	TypeBool   = "bool"
	TypeString = "string"
)

type commandArgument struct {
	label       string
	description string
	position    int
	valueType   string
}

type commandOption struct {
	label       string
	description string
	letter      rune
	name        string
	valueType   *string // if the option does not need a value it shall have a valueType of nil
}

type CommandInput interface {
	GetArgument(commandArgument) any
	GetOption(commandOption) any
}

type commandInput struct {
	arguments map[commandArgument]any
	options   map[commandOption]any
}

func (c *commandInput) GetArgument(arg commandArgument) any {
	return c.arguments[arg]
}

func (c *commandInput) GetOption(opt commandOption) any {
	return c.options[opt]
}

type Command interface {
	AddArgument(int, string) Command
	AddOption(rune, string, string) Command
	Parse([]string) (CommandInput, Error)
}

type command struct {
	Name      string
	Arguments []commandArgument
	Options   []commandOption
}

func NewCommand(name string) Command {
	return &command{Name: name}
}

func (c *command) AddArgument(position int, valueType string) Command {
	c.Arguments = append(c.Arguments, commandArgument{
		position:  position,
		valueType: valueType,
	})
	return c
}

func (c *command) AddOption(letter rune, name string, valueType string) Command {
	c.Options = append(c.Options, commandOption{
		letter:    letter,
		name:      name,
		valueType: &valueType,
	})
	return c
}

func (c *command) Parse(input []string) (CommandInput, Error) {
	inputLength := len(input)

	if inputLength == 0 {
		return nil, &InvalidCommandUsageError{command: c.Name}
	}

	inputArgs := make(map[commandArgument]any)
	inputOpts := make(map[commandOption]any)

	// Parse arguments
	nbrArguments := len(c.Arguments)
	if inputLength < nbrArguments {
		return nil, &InvalidCommandUsageError{command: c.Name}
	}
	for _, arg := range c.Arguments {
		value, err := ParseValue(arg.valueType, input[arg.position])
		if err != nil {
			return nil, &InvalidCommandUsageError{command: c.Name}
		}
		inputArgs[arg] = value
	}

	// Parse options
	for _, opt := range c.Options {
		index := slices.Index(input, string(opt.letter))
		if index == -1 {
			index = slices.Index(input, opt.name)
		}
		if index != -1 {
			if opt.valueType == nil {
				inputOpts[opt] = true
			} else {
				if index+1 >= inputLength {
					return nil, &InvalidCommandUsageError{command: c.Name}
				}
				value, err := ParseValue(*opt.valueType, input[index+1])
				if err != nil {
					return nil, &InvalidCommandUsageError{command: c.Name}
				}
				inputOpts[opt] = value
			}
		}
	}

	return &commandInput{
		arguments: inputArgs,
		options:   inputOpts,
	}, nil
}
