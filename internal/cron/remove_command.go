package cron

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/internal/database"
)

type removeCommand struct {
	scheduler scheduler
}

func (cmd *removeCommand) Name() string {
	return "remove"
}

func (cmd *removeCommand) HelpCommand() string {
	return "remove <name>"
}

func (cmd *removeCommand) Description() string {
	return "Remove a schedule of the specified name"
}

func (cmd *removeCommand) Execute(ctx context.Context, params []string, channel string) (string, error) {
	params = params[1:]
	if len(params) != 1 {
		return "", CommandSyntaxError
	}

	db, ok := database.FromContext(ctx)
	if !ok {
		return "", errors.New("failed to get database from context")
	}

	name := params[0]

	err := db.DeleteScheduleByName(ctx, name)
	if err != nil {
		if err == database.ErrNotFound {
			return fmt.Sprintf("%s does not exist.", name), nil
		}
		return "", fmt.Errorf("failed to remove a schedule %s: %w", name, err)
	}

	cmd.scheduler.removeSchedule(ctx, name)

	return fmt.Sprintf("Success to remove a schedule : %s", name), nil
}
