package main

import (
	"fmt"
)

type InvalidCommandError struct {
	command string
}

func (e *InvalidCommandError) Error() string {
	return fmt.Sprintf("Invalid command: %s", e.command)
}

type UnexpectedError struct {
	message string
	err     error
}

func (e *UnexpectedError) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err)
}
