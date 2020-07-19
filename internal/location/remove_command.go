package location

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/v2/internal/database"
)

type removeCommand struct{}

func (cmd *removeCommand) Name() string {
	return "remove"
}

func (cmd *removeCommand) HelpCommand() string {
	return "remove <name>"
}

func (cmd *removeCommand) Description() string {
	return "Remove a location of the specified name"
}

func (cmd *removeCommand) Execute(ctx context.Context, params []string) (string, error) {
	params = params[1:]
	if len(params) != 1 {
		return "", CommandSyntaxError
	}

	db, ok := database.FromContext(ctx)
	if !ok {
		return "", errors.New("failed to get database from context")
	}

	name := params[0]

	err := db.DeleteLocationByName(ctx, name)
	if err != nil {
		if err == database.ErrNotFound {
			return fmt.Sprintf("%s does not exist.", name), nil
		}
		return "", fmt.Errorf("failed to remove a location %s: %w", name, err)
	}

	return fmt.Sprintf("Success to remove a location : %s", name), nil
}
