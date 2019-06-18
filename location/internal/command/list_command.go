package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/kechako/gopher-bot/location/internal/data"
	"golang.org/x/xerrors"
)

type listCommand struct{}

func (cmd *listCommand) Name() string {
	return "list"
}

func (cmd *listCommand) HelpCommand() string {
	return "list"
}

func (cmd *listCommand) Description() string {
	return "List locations."
}

func (cmd *listCommand) Execute(ctx context.Context, params []string) (string, error) {
	if len(params) != 1 {
		return "", CommandSyntaxError
	}

	locs, err := data.GetLocations(ctx)
	if err != nil {
		return "", xerrors.Errorf("failed to get locations: %w", err)
	}

	var msg strings.Builder
	for i, loc := range locs {
		if i > 0 {
			msg.WriteString("\n")
		}
		msg.WriteString(fmt.Sprintf("%s [%f, %f]", loc.Name, loc.Latitude, loc.Longitude))
	}

	return msg.String(), nil
}
