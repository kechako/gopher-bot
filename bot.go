package bot

import (
	"context"
	"strings"
	"time"

	"github.com/kechako/gopher-bot/plugin"
	"github.com/kechako/gopher-bot/service"
)

// Bot represents a bot.
type Bot struct {
	service service.Service
	plugins []plugin.Plugin
}

// New returns a new *Bot.
func New(s service.Service) *Bot {
	return &Bot{
		service: s,
	}
}

// AddPlugin adds a plugin to the bot.
func (b *Bot) AddPlugin(p plugin.Plugin) {
	b.plugins = append(b.plugins, p)
}

// Run runs the bot.
func (b *Bot) Run(ctx context.Context) error {
	ch, err := b.service.Start(ctx)
	if err != nil {
		return err
	}

	b.loop(ctx, ch)

	return nil
}

func (b *Bot) loop(ctx context.Context, ch <-chan *service.Event) {
	defer b.service.Close()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case event := <-ch:
			switch event.Type {
			case service.ConnectedEvent:
				if hello := event.GetHello(); hello != nil {
					b.hello(ctx, hello)
				}
			case service.MessageEvent:
				if msg := event.GetMessage(); msg != nil {
					b.doAction(ctx, msg)
				}
			}
		}
	}
}

func (b *Bot) hello(ctx context.Context, hello plugin.Hello) {
	for _, p := range b.plugins {
		callPluginHello(ctx, p, hello)
	}
}

func callPluginHello(ctx context.Context, plugin plugin.Plugin, hello plugin.Hello) {
	defer func() {
		if err := recover(); err != nil {
			// TODO: output error log
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	plugin.Hello(ctx, hello)
}

func (b *Bot) doAction(ctx context.Context, msg plugin.Message) {
	if msg.MentionTo(b.service.UserID()) && strings.Contains(msg.Text(), "help") {
		b.postHelp(ctx, msg.ChannelID())
		return
	}

	for _, p := range b.plugins {
		callPluginDoAction(ctx, p, msg)
	}
}

func callPluginDoAction(ctx context.Context, plugin plugin.Plugin, msg plugin.Message) {
	defer func() {
		if err := recover(); err != nil {
			// TODO: output error log
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	plugin.DoAction(ctx, msg)
}

func (b *Bot) postHelp(ctx context.Context, channelID string) {
	var doc strings.Builder

	for i, p := range b.plugins {
		h := callPluginHelp(ctx, p)
		if h == nil {
			continue
		}
		if i > 0 {
			doc.WriteString("\n")
		}

		doc.WriteString(h.String())
	}

	escaped := b.service.EscapeHelp(doc.String())
	b.service.Post(channelID, escaped)
}

func callPluginHelp(ctx context.Context, plugin plugin.Plugin) *plugin.Help {
	defer func() {
		if err := recover(); err != nil {
			// TODO: output error log
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return plugin.Help(ctx)
}
