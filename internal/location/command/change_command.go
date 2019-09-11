package command

import (
	"context"
	"fmt"

	"github.com/kechako/gopher-bot/internal/location/data"
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

	if err := data.UpdateLocation(ctx, loc); err != nil {
		if err == data.ErrKeyNotFound {
			return fmt.Sprintf("%s does not exist.", loc.Name), nil
		}
		return "", fmt.Errorf("failed to update a location %s: %w", loc.Name, err)
	}

	return fmt.Sprintf("Success to change a location : %s [%f, %f]", loc.Name, loc.Latitude, loc.Longitude), nil
}
