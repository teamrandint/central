package commands

import (
	"time"
)

type Commands interface {
	HasTimeElapsed() bool
	SetCreationTime()
	GetCreationTime()
	SetCommandName()
	GetCommandName() string
}

type Command struct {
	commandName string
	creationTime time.Time
	user string
	params []string
}

// Constructor for creating a new command
func NewCommand(name string, userId string, cmdParams []string) *Command {
	cmd := new(Command)
	cmd.commandName = name
	cmd.user = userId
	cmd.params = cmdParams
	cmd.SetCreationTime()
	return cmd
}

// Determines if the command has been executed in the last 60 seconds or not
// Returns true if time since creation is 60 or greater, false otherwise.
func (cmd *Command) HasTimeElapsed() bool {
	currentTime := time.Now()
	elapsed := currentTime.Sub(cmd.creationTime)

	if elapsed.Minutes() >= 1 {
		return true
	}

	return false
}

func (cmd *Command) SetCreationTime() {
	cmd.creationTime = time.Now()
}

func (cmd *Command) CreationTime() time.Time {
	return cmd.creationTime
}

func (cmd *Command) SetCommandName(name string) {
	cmd.commandName = name
}

func (cmd *Command) CommandName() string {
	return cmd.commandName
}