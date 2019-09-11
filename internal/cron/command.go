package cron

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/internal/cron/data"
	"github.com/kechako/gopher-bot/plugin"
	cron "github.com/robfig/cron/v3"
)

type scheduler interface {
	addSchedule(ctx context.Context, s *data.Schedule) error
	resetScheduler(ctx context.Context) error
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
	cron         *cron.Cron

	bot Bot
}

var _ scheduler = (*Cron)(nil)

func New(bot Bot) *Cron {
	c := &Cron{
		commanderMap: make(map[string]Commander),
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
	return c.resetScheduler(ctx)
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

func (c *Cron) addSchedule(ctx context.Context, s *data.Schedule) error {
	_, err := c.cron.AddFunc(s.Fields, cron.FuncJob(func() {
		c.bot.ProcessCommand(s.Channel, s.Command)
	}))
	if err != nil {
		return err
	}

	return nil
}

func (c *Cron) resetScheduler(ctx context.Context) error {
	if c.cron != nil {
		c.cron.Stop()
	}

	schedules, err := data.GetSchedules(ctx)
	if err != nil {
		return err
	}

	c.cron = cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)))
	for _, s := range schedules {
		c.addSchedule(ctx, s)
	}
	c.cron.Start()

	return nil
}
