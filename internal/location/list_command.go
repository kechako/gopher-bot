package location

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kechako/gopher-bot/v2/internal/database"
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

	db, ok := database.FromContext(ctx)
	if !ok {
		return "", errors.New("failed to get database from context")
	}

	locs, err := db.SearchLocations(ctx)
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
