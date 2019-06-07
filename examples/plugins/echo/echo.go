package echo

import (
	"github.com/kechako/gopher-bot/plugin"
)

type echo struct {
	bot plugin.Bot
}

var _ plugin.Plugin = (*echo)(nil)

func New() plugin.Plugin {
	return &echo{}
}

func (e *echo) Hello(hello plugin.Hello) {
	e.bot = hello.Bot()
}

func (e *echo) DoAction(msg plugin.Message) {
	msg.Post(msg.Text())
}

func (e *echo) Help() *plugin.Help {
	return &plugin.Help{
		Name:        "echo",
		Description: "echo plugin posts echo message",
	}
}
