package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/kechako/gopher-bot/location/internal/data"
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
		return "", fmt.Errorf("failed to get locations: %w", err)
	}

	if len(locs) == 0 {
		return "Location list is empty.", nil
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
