package cron

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kechako/gopher-bot/internal/database"
)

type listCommand struct {
	bot Bot
}

func (cmd *listCommand) Name() string {
	return "list"
}

func (cmd *listCommand) HelpCommand() string {
	return "list"
}

func (cmd *listCommand) Description() string {
	return "List schedules."
}

func (cmd *listCommand) Execute(ctx context.Context, params []string, channel string) (string, error) {
	if len(params) != 1 {
		return "", CommandSyntaxError
	}

	db, ok := database.FromContext(ctx)
	if !ok {
		return "", errors.New("failed to get database from context")
	}

	sches, err := db.SearchSchedules(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get schedules: %w", err)
	}

	if len(sches) == 0 {
		return "Schedule list is empty.", nil
	}

	var msg strings.Builder
	for i, sch := range sches {
		if i > 0 {
			msg.WriteString("\n")
		}
		msg.WriteString(fmt.Sprintf("%s : %s %s [%s]", sch.Name, sch.Fields, sch.Command, cmd.bot.ChannelName(sch.Channel)))
	}

	return msg.String(), nil
}
