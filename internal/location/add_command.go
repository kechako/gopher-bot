package location

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/kechako/gopher-bot/v2/internal/database"
)

type addCommand struct{}

func (cmd *addCommand) Name() string {
	return "add"
}

func (cmd *addCommand) HelpCommand() string {
	return "add <name> <latitude> <longitude>"
}

func (cmd *addCommand) Description() string {
	return "Add a new location with specified name."
}

func (cmd *addCommand) Execute(ctx context.Context, params []string) (string, error) {
	params = params[1:]
	loc, err := makeLocation(params)
	if err != nil {
		return "", err
	}

	db, ok := database.FromContext(ctx)
	if !ok {
		return "", errors.New("failed to get database from context")
	}

	if err := db.SaveLocation(ctx, loc); err != nil {
		if err == database.ErrDuplicated {
			return fmt.Sprintf("%s already exists", loc.Name), nil
		}

		return "", fmt.Errorf("failed to add a new location %s: %w", loc.Name, err)
	}

	return fmt.Sprintf("Success to add a new location : %s [%f, %f]", loc.Name, loc.Latitude, loc.Longitude), nil
}

func makeLocation(params []string) (*database.Location, error) {
	if len(params) != 3 {
		return nil, ErrInvalidSyntax
	}

	name := params[0]

	lat, err := strconv.ParseFloat(params[1], 32)
	if err != nil {
		return nil, ErrInvalidSyntax
	}

	lon, err := strconv.ParseFloat(params[2], 32)
	if err != nil {
		return nil, ErrInvalidSyntax
	}

	return &database.Location{
		Name:      name,
		Latitude:  float32(lat),
		Longitude: float32(lon),
	}, nil
}
