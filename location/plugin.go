package location

import (
	"context"
	"strings"

	"github.com/kechako/gopher-bot/internal/location/command"
	"github.com/kechako/gopher-bot/plugin"
)

const commandName = "loc"

type locationPlugin struct {
	cmd *command.Command
}

var _ plugin.Plugin = (*locationPlugin)(nil)

// NewPlugin returns a new plugin.Plugin that manages locations.
func NewPlugin() plugin.Plugin {
	return &locationPlugin{
		cmd: command.New(),
	}
}

func (p *locationPlugin) Hello(ctx context.Context, hello plugin.Hello) {
}

func (p *locationPlugin) DoAction(ctx context.Context, msg plugin.Message) {
	params := strings.Fields(msg.Text())
	if len(params) == 0 || params[0] != commandName {
		return
	}

	retMsg, err := p.cmd.Execute(ctx, params[1:])
	if err != nil {
		if err == command.CommandSyntaxError {
			msg.PostHelp(p.Help(ctx))
			return
		}
		// TODO: output error log
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
