package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/plugin"
)

var (
	CommandSyntaxError = errors.New("CommandSyntaxError")
)

type Commander interface {
	Name() string
	HelpCommand() string
	Description() string
	Execute(ctx context.Context, params []string) (string, error)
}

type Command struct {
	commanders   []Commander
	commanderMap map[string]Commander
}

func New() *Command {
	cmd := &Command{
		commanders: []Commander{
			&addCommand{},
			&listCommand{},
			&removeCommand{},
			&changeCommand{},
			&helpCommand{},
		},
		commanderMap: make(map[string]Commander),
	}
	cmd.init()
	return cmd
}

func (cmd *Command) init() {
	for _, cmdr := range cmd.commanders {
		name := cmdr.Name()
		cmd.commanderMap[name] = cmdr
	}
}

func (cmd *Command) Execute(ctx context.Context, params []string) (string, error) {
	var cmdName string
	if len(params) > 0 {
		cmdName = string(params[0])
	}

	commander, ok := cmd.commanderMap[cmdName]
	if !ok {
		return "", CommandSyntaxError
	}

	return commander.Execute(ctx, params)
}

func (cmd *Command) HelpCommands(name string) []*plugin.Command {
	var commands []*plugin.Command

	for _, cmdr := range cmd.commanders {
		command := &plugin.Command{
			Command:     fmt.Sprintf("%s %s", name, cmdr.HelpCommand()),
			Description: cmdr.Description(),
		}
		commands = append(commands, command)
	}

	return commands
}
