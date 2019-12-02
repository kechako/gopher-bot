// Package echo is a plugin to echo messages.
package echo

import (
	"context"
	"github.com/kechako/gopher-bot/plugin"
)

type echo struct {
	bot plugin.Bot
}

var _ plugin.Plugin = (*echo)(nil)

func New() plugin.Plugin {
	return &echo{}
}

func (e *echo) Hello(ctx context.Context, hello plugin.Hello) {
	e.bot = hello.Bot()
}

func (e *echo) DoAction(ctx context.Context, msg plugin.Message) {
	msg.Post(msg.Text())
}

func (e *echo) Help(ctx context.Context) *plugin.Help {
	return &plugin.Help{
		Name:        "echo",
		Description: "echo plugin posts echo message",
	}
}
