// Package cron provides commands to manage cron schedules.
package cron

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/v2/internal/database"
	"github.com/kechako/gopher-bot/v2/plugin"
	cron "github.com/robfig/cron/v3"
)

type scheduler interface {
	addSchedule(ctx context.Context, s *database.Schedule) error
	removeSchedule(ctx context.Context, name string)
}

var (
	CommandSyntaxError = errors.New("CommandSyntaxError")
)

type Bot interface {
	ProcessCommand(channelID, command string)
	ChannelName(channelID string) string
}

type Commander interface {
	Name() string
	HelpCommand() string
	Description() string
	Execute(ctx context.Context, params []string, channel string) (string, error)
}

type CommandFunc func(channelID string, command string)

type Cron struct {
	commanders   []Commander
	commanderMap map[string]Commander

	cron    *cron.Cron
	entries map[string]cron.EntryID

	bot Bot
}

var _ scheduler = (*Cron)(nil)

func New(bot Bot) *Cron {
	c := &Cron{
		commanderMap: make(map[string]Commander),
		cron:         cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow))),
		entries:      make(map[string]cron.EntryID),
		bot:          bot,
	}
	c.commanders = []Commander{
		&addCommand{
			scheduler: c,
			bot:       bot,
		},
		&listCommand{
			bot: bot,
		},
		&removeCommand{
			scheduler: c,
		},
		&helpCommand{},
	}
	c.init()

	return c
}

func (c *Cron) init() {
	for _, cmdr := range c.commanders {
		name := cmdr.Name()
		c.commanderMap[name] = cmdr
	}
}

func (c *Cron) Start(ctx context.Context) error {
	db, ok := database.FromContext(ctx)
	if !ok {
		return errors.New("failed to get database from context")
	}

	schedules, err := db.SearchSchedules(ctx)
	if err != nil {
		return err
	}

	for _, s := range schedules {
		c.addSchedule(ctx, s)
	}

	c.cron.Start()

	return nil
}

func (c *Cron) Close() error {
	c.cron.Stop()
	return nil
}

func (c *Cron) Execute(ctx context.Context, params []string, channel string) (string, error) {
	var cmdName string
	if len(params) > 0 {
		cmdName = string(params[0])
	}

	commander, ok := c.commanderMap[cmdName]
	if !ok {
		return "", CommandSyntaxError
	}

	return commander.Execute(ctx, params, channel)
}

func (c *Cron) HelpCommands(name string) []*plugin.Command {
	var commands []*plugin.Command

	for _, cmdr := range c.commanders {
		command := &plugin.Command{
			Command:     fmt.Sprintf("%s %s", name, cmdr.HelpCommand()),
			Description: cmdr.Description(),
		}
		commands = append(commands, command)
	}

	return commands
}

func (c *Cron) addSchedule(ctx context.Context, s *database.Schedule) error {
	id, err := c.cron.AddFunc(s.Fields, cron.FuncJob(func() {
		c.bot.ProcessCommand(s.Channel, s.Command)
	}))
	if err != nil {
		return err
	}

	c.entries[s.Name] = id

	return nil
}

func (c *Cron) removeSchedule(ctx context.Context, name string) {
	id, ok := c.entries[name]
	if !ok {
		return
	}

	c.cron.Remove(id)
}
