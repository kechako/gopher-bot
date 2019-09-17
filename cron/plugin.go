package cron

import (
	"context"
	"strings"

	"github.com/kechako/gopher-bot/internal/cron"
	"github.com/kechako/gopher-bot/logger"
	"github.com/kechako/gopher-bot/plugin"
)

const commandName = "cron"

type cronPlugin struct {
	bot  plugin.Bot
	cron *cron.Cron
}

var _ plugin.Plugin = (*cronPlugin)(nil)

// New returns a new plugin.Plugin that manages crons.
func New() plugin.Plugin {
	return &cronPlugin{}
}

func (p *cronPlugin) Close() error {
	return p.cron.Close()
}

func (p *cronPlugin) Hello(ctx context.Context, hello plugin.Hello) {
	p.bot = hello.Bot()
	p.cron = cron.New(&cronBot{plugin: p})

	if err := p.cron.Start(ctx); err != nil {
		logger.FromContext(ctx).Error("failed to start cron: ", err)
	}
}

func (p *cronPlugin) DoAction(ctx context.Context, msg plugin.Message) {
	params := strings.Fields(msg.Text())
	if len(params) == 0 || params[0] != commandName {
		return
	}

	retMsg, err := p.cron.Execute(ctx, params[1:], msg.ChannelID())
	if err != nil {
		if err == cron.CommandSyntaxError {
			msg.PostHelp(p.Help(ctx))
			return
		}

		logger.FromContext(ctx).Error("cron: ", err)
		return
	}

	msg.Post(retMsg)
}

func (p *cronPlugin) Help(ctx context.Context) *plugin.Help {
	return &plugin.Help{
		Name:        "cron",
		Description: "Management command schedules.",
		Commands:    p.cron.HelpCommands(commandName),
	}
}

type cronBot struct {
	plugin *cronPlugin
}

var _ cron.Bot = (*cronBot)(nil)

func (bot *cronBot) ProcessCommand(channelID, command string) {
	bot.plugin.bot.ProcessCommand(channelID, command)
}

func (bot *cronBot) ChannelName(channelID string) string {
	ch := bot.plugin.bot.Channel(channelID)
	if ch == nil {
		return ""
	}

	return ch.Name()
}
