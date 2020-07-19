package location

import (
	"context"
	"errors"
	"fmt"

	"github.com/kechako/gopher-bot/internal/database"
)

type changeCommand struct{}

func (cmd *changeCommand) Name() string {
	return "change"
}

func (cmd *changeCommand) HelpCommand() string {
	return "change <name> <latitude> <longitude>"
}

func (cmd *changeCommand) Description() string {
	return "Change a location of the specified name."
}

func (cmd *changeCommand) Execute(ctx context.Context, params []string) (string, error) {
	params = params[1:]
	loc, err := makeLocation(params)
	if err != nil {
		return "", err
	}

	db, ok := database.FromContext(ctx)
	if !ok {
		return "", errors.New("failed to get database from context")
	}

	oldloc, err := db.FindLocationByName(ctx, loc.Name)
	if err != nil {
		if err == database.ErrNotFound {
			return fmt.Sprintf("%s does not exist.", loc.Name), nil
		}
		return "", fmt.Errorf("failed to update a location %s: %w", loc.Name, err)
	}

	loc.ID = oldloc.ID

	if err := db.SaveLocation(ctx, loc); err != nil {
		if err == database.ErrNotFound {
			return fmt.Sprintf("%s does not exist.", loc.Name), nil
		}
		return "", fmt.Errorf("failed to update a location %s: %w", loc.Name, err)
	}

	return fmt.Sprintf("Success to change a location : %s [%f, %f]", loc.Name, loc.Latitude, loc.Longitude), nil
}
