package main

import (
	"fmt"
	"slices"
	"strconv"
)

const (
	TypeInt    = "int"
	TypeFloat  = "float"
	TypeBool   = "bool"
	TypeString = "string"
)

// ParseValue attempts to parse a string into the specified type.
func ParseValue(valueType string, value string) (interface{}, error) {
	switch valueType {
	case TypeInt:
		return strconv.Atoi(value) // Parse string to int
	case TypeFloat:
		return strconv.ParseFloat(value, 64) // Parse string to float64
	case TypeBool:
		return strconv.ParseBool(value) // Parse string to bool
	case TypeString:
		return value, nil // No parsing needed for strings
	default:
		return nil, fmt.Errorf("unsupported type: %s", valueType)
	}
}

type CommandArgument struct {
	position  int
	valueType string
}

// if the option does not need a value it shall have a valueType of nil
type CommandOption struct {
	letter    rune
	name      string
	valueType *string
}

type CommandInput struct {
	arguments map[CommandArgument]any
	options   map[CommandOption]any
}

type commandOutput struct {
	output string
}

func (c *commandOutput) String() string {
	return c.output
}

type CommandOutput interface {
	String() string
}

type Command struct {
	Name      string
	Arguments []CommandArgument
	Options   []CommandOption
	handler   func(args CommandInput) (CommandOutput, Error)
}

func NewCommand(name string) *Command {
	return &Command{Name: name}
}

func (command *Command) AddArgument(position int, valueType string) *Command {
	command.Arguments = append(command.Arguments, CommandArgument{
		position:  position,
		valueType: valueType,
	})
	return command
}

func (command *Command) AddOption(letter rune, name string, valueType string) *Command {
	command.Options = append(command.Options, CommandOption{
		letter:    letter,
		name:      name,
		valueType: &valueType,
	})
	return command
}

func (command *Command) SetHandler(handler func(args CommandInput) (CommandOutput, Error)) *Command {
	command.handler = handler
	return command
}

func (command *Command) Parse(input []string) (*CommandInput, Error) {
	inputLength := len(input)

	if inputLength == 0 {
		return nil, &InvalidCommandUsageError{command: command.Name}
	}

	inputArgs := make(map[CommandArgument]any)
	inputOpts := make(map[CommandOption]any)

	// Parse arguments
	nbrArguments := len(command.Arguments)
	if inputLength < nbrArguments {
		return nil, &InvalidCommandUsageError{command: command.Name}
	}
	for _, arg := range command.Arguments {
		value, err := ParseValue(arg.valueType, input[arg.position])
		if err != nil {
			return nil, &InvalidCommandUsageError{command: command.Name}
		}
		inputArgs[arg] = value
	}

	// Parse options
	for _, opt := range command.Options {
		index := slices.Index(input, string(opt.letter))
		if index == -1 {
			index = slices.Index(input, opt.name)
		}
		if index != -1 {
			if opt.valueType == nil {
				inputOpts[opt] = true
			} else {
				if index+1 >= inputLength {
					return nil, &InvalidCommandUsageError{command: command.Name}
				}
				value, err := ParseValue(*opt.valueType, input[index+1])
				if err != nil {
					return nil, &InvalidCommandUsageError{command: command.Name}
				}
				inputOpts[opt] = value
			}
		}
	}

	return &CommandInput{
		arguments: inputArgs,
		options:   inputOpts,
	}, nil
}

func (command *Command) Run(input CommandInput) (CommandOutput, Error) {
	return &commandOutput{output: "OK"}, nil
}

func InitCommands() map[string]*Command {
	getCommand := NewCommand("GET").
		AddArgument(0, "string").
		AddOption('e', "expires-in", "int")

	return map[string]*Command{
		"GET": getCommand,
	}
}
