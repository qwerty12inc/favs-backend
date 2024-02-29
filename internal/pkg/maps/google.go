package maps

import (
	"fmt"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"regexp"
	"strconv"
)

type LocationLinkResolver interface {
	ResolveLink(link string) (models.Coordinates, error)
}

type LocationLinkResolverImpl struct {
}

var coordinatesRegexp = `@(-?\d+\.\d+),(-?\d+\.\d+)`

func (l LocationLinkResolverImpl) ResolveLink(link string) (models.Coordinates, error) {
	r, err := regexp.Compile(coordinatesRegexp)
	if err != nil {
		return models.Coordinates{}, err
	}

	matches := r.FindStringSubmatch(link)
	if len(matches) == 3 {
		latStr := matches[1]
		lonStr := matches[2]
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return models.Coordinates{}, err
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return models.Coordinates{}, err
		}

		return models.Coordinates{
			Latitude:  lat,
			Longitude: lon,
		}, nil
	}

	return models.Coordinates{}, fmt.Errorf("link %s does not contain coordinates", link)
}
