package command

import (
	"context"
	"fmt"

	"github.com/kechako/gopher-bot/location/internal/data"
)

type removeCommand struct{}

func (cmd *removeCommand) Name() string {
	return "remove"
}

func (cmd *removeCommand) HelpCommand() string {
	return "remove <name> <latitude> <longitude>"
}

func (cmd *removeCommand) Description() string {
	return "Remove a location of the specified name"
}

func (cmd *removeCommand) Execute(ctx context.Context, params []string) (string, error) {
	params = params[1:]
	if len(params) != 1 {
		return "", CommandSyntaxError
	}

	name := params[0]

	err := data.RemoveLocation(ctx, name)
	if err != nil {
		if err == data.ErrKeyNotFound {
			return fmt.Sprintf("%s does not exist.", name), nil
		}
		return "", fmt.Errorf("failed to remove a location %s: %w", name, err)
	}

	return fmt.Sprintf("Success to remove a location : %s", name), nil
}
