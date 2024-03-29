package cron

import (
	"context"
)

type helpCommand struct{}

func (cmd *helpCommand) Name() string {
	return "help"
}

func (cmd *helpCommand) HelpCommand() string {
	return "help"
}

func (cmd *helpCommand) Description() string {
	return "Show this help message."
}

func (cmd *helpCommand) Execute(ctx context.Context, params []string, channel string) (string, error) {
	return "", ErrInvalidSyntax
}
