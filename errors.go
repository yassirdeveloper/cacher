package main

import (
	"fmt"
)

type Error interface {
	Error() string
	Display() string
}

type InvalidCommandError struct {
	command string
}

func (e *InvalidCommandError) Error() string {
	return fmt.Sprintf("Invalid command: %s", e.command)
}

func (e *InvalidCommandError) Display() string {
	return fmt.Sprintf("Invalid command: %s", e.command)
}

type InvalidCommandUsageError struct {
	command string
}

func (e *InvalidCommandUsageError) Error() string {
	return fmt.Sprintf("Invalid usage of command: %s", e.command)
}

func (e *InvalidCommandUsageError) Display() string {
	return fmt.Sprintf("Invalid usage of command: %s", e.command)
}

type CommandNotExecutableError struct {
	command string
}

func (e *CommandNotExecutableError) Error() string {
	return fmt.Sprintf("Command not executable: %s", e.command)
}

func (e *CommandNotExecutableError) Display() string {
	return fmt.Sprintf("Command not executable: %s", e.command)
}

type CommandError struct {
	message string
}

func (e *CommandError) Error() string {
	return e.message
}

func (e *CommandError) Display() string {
	return e.message
}

type UnexpectedError struct {
	message string
	err     error
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err)
}

func (e *UnexpectedError) Display() string {
	return "Something went wrong!"
}

type SetupError struct {
	message string
}

func (e *SetupError) Error() string {
	return e.message
}

func (e *SetupError) Display() string {
	return e.message
}
