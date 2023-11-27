// Package location is a plugin to manage locations.
package location

import (
	"context"
	"log/slog"
	"strings"

	"github.com/kechako/gopher-bot/v2/internal/location"
	"github.com/kechako/gopher-bot/v2/plugin"
)

const commandName = "loc"

type locationPlugin struct {
	cmd *location.Command
	l   *slog.Logger
}

var _ plugin.Plugin = (*locationPlugin)(nil)

// NewPlugin returns a new plugin.Plugin that manages locations.
func NewPlugin() plugin.Plugin {
	return &locationPlugin{
		cmd: location.New(),
	}
}

func (p *locationPlugin) Hello(ctx context.Context, hello plugin.Hello) {
	p.l = hello.Bot().Logger().With(slog.String("plugin", "location"))
}

func (p *locationPlugin) DoAction(ctx context.Context, msg plugin.Message) {
	params := strings.Fields(msg.Text())
	if len(params) == 0 || params[0] != commandName {
		return
	}

	retMsg, err := p.cmd.Execute(ctx, params[1:])
	if err != nil {
		if err == location.CommandSyntaxError {
			msg.PostHelp(p.Help(ctx))
			return
		}

		p.l.Error("failed to do plugin action", slog.Any("err", err))
		return
	}

	msg.Post(retMsg)
}

func (p *locationPlugin) Help(ctx context.Context) *plugin.Help {
	return &plugin.Help{
		Name:        "location",
		Description: "Management location. Locations are used by each plugin.",
		Commands:    p.cmd.HelpCommands(commandName),
	}
}
