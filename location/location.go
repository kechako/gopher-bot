package location

import (
	"context"

	"github.com/kechako/gopher-bot/location/internal/data"
)

// GetLocation returns the latitude and longitude of the specified name.
func GetLocation(ctx context.Context, name string) (lat, lon float32, err error) {
	loc, err := data.GetLocation(ctx, name)
	if err != nil {
		return
	}

	lat = loc.Latitude
	lon = loc.Longitude

	return
}
