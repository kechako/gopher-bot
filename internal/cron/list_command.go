package cron

import (
	"context"
	"fmt"
	"strings"

	"github.com/kechako/gopher-bot/internal/cron/data"
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

	sches, err := data.GetSchedules(ctx)
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
