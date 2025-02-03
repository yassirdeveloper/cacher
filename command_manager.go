package main

type CommandManager interface {
	Get(string) (Command, Error)
	AddCommand(string, Command) CommandManager
}

type commandManager struct {
	commands map[string]Command
}

func NewCommandManager() CommandManager {
	command_manager := &commandManager{}
	command_manager.AddCommand("GET", GetCommand)
	return command_manager
}

func (c *commandManager) AddCommand(commandName string, command Command) CommandManager {
	c.commands[commandName] = command
	return c
}

func (c *commandManager) Get(commandName string) (Command, Error) {
	command := c.commands[commandName]
	if command == nil {
		return nil, &InvalidCommandError{command: commandName}
	}
	return command, nil
}
