package cron

import (
	"context"
	"fmt"

	"github.com/kechako/gopher-bot/internal/cron/data"
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

	name := params[0]

	err := data.RemoveSchedule(ctx, name)
	if err != nil {
		if err == data.ErrKeyNotFound {
			return fmt.Sprintf("%s does not exist.", name), nil
		}
		return "", fmt.Errorf("failed to remove a schedule %s: %w", name, err)
	}

	if err := cmd.scheduler.resetScheduler(ctx); err != nil {
		return "", fmt.Errorf("failed to reset scheduler: %w", err)
	}

	return fmt.Sprintf("Success to remove a schedule : %s", name), nil
}
