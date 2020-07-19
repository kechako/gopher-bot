package location

import (
	"context"
	"errors"

	"github.com/kechako/gopher-bot/internal/database"
)

var (
	// ErrLocationNotFound is the error used for the location not found.
	ErrLocationNotFound = errors.New("the location not found")
)

// GetLocation returns the latitude and longitude of the specified name.
func GetLocation(ctx context.Context, name string) (lat, lon float32, err error) {
	db, ok := database.FromContext(ctx)
	if !ok {
		return 0, 0, errors.New("failed to get database from context")
	}

	loc, err := db.FindLocationByName(ctx, name)
	if err != nil {
		if err == database.ErrNotFound {
			return 0, 0, ErrLocationNotFound
		}
		return
	}

	lat = loc.Latitude
	lon = loc.Longitude

	return
}
