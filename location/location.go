package location

import (
	"context"
	"errors"

	"github.com/kechako/gopher-bot/location/internal/data"
)

var (
	ErrLocationNotFound = errors.New("the location not found")
)

// GetLocation returns the latitude and longitude of the specified name.
func GetLocation(ctx context.Context, name string) (lat, lon float32, err error) {
	loc, err := data.GetLocation(ctx, name)
	if err != nil {
		if err == data.ErrKeyNotFound {
			return 0, 0, ErrLocationNotFound
		}
		return
	}

	lat = loc.Latitude
	lon = loc.Longitude

	return
}
