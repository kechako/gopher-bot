package cron

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kechako/gopher-bot/v2/internal/database"
)

type addCommand struct {
	scheduler scheduler
	bot       Bot
}

func (cmd *addCommand) Name() string {
	return "add"
}

func (cmd *addCommand) HelpCommand() string {
	return "add <name> <schedule> <command>"
}

func (cmd *addCommand) Description() string {
	return "Add a new schedule with specified name."
}

func (cmd *addCommand) Execute(ctx context.Context, params []string, channel string) (string, error) {
	params = params[1:]
	sch, err := makeSchedule(params, channel)
	if err != nil {
		return "", err
	}

	db, ok := database.FromContext(ctx)
	if !ok {
		return "", errors.New("failed to get database from context")
	}

	if err := cmd.scheduler.addSchedule(ctx, sch); err != nil {
		return "", ErrInvalidSyntax
	}

	if err := db.SaveSchedule(ctx, sch); err != nil {
		if err == database.ErrDuplicated {
			return fmt.Sprintf("%s already exists", sch.Name), nil
		}

		return "", fmt.Errorf("failed to add a new schedule %s: %w", sch.Name, err)
	}

	return fmt.Sprintf("Success to add a new schedule : %s [%s, %s, %s]", sch.Name, sch.Fields, sch.Command, cmd.bot.ChannelName(sch.Channel)), nil
}

func makeSchedule(params []string, channel string) (*database.Schedule, error) {
	if len(params) < 7 {
		return nil, ErrInvalidSyntax
	}

	name := params[0]
	fields := strings.Join(params[1:6], " ")
	command := strings.Join(params[6:], " ")

	return &database.Schedule{
		Name:    name,
		Fields:  fields,
		Command: command,
		Channel: channel,
	}, nil
}
