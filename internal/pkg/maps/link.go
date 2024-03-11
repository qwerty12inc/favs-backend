package maps

import (
	"fmt"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type LocationLinkResolver interface {
	// ResolveLink resolves a location link and returns coordinates
	ResolveLink(link string) (models.Coordinates, error)
}

type LocationLinkResolverImpl struct {
}

var coordinatesRegexp = `@(-?\d+\.\d+),(-?\d+\.\d+)`
var googleMapsRegexp = `https:\/\/maps\.app\.goo\.gl\/[A-Za-z0-9]+`

func (l LocationLinkResolverImpl) ResolveLink(link string) (models.Coordinates, error) {
	r, err := regexp.Compile(googleMapsRegexp)
	if err != nil {
		return models.Coordinates{}, err
	}

	if r.MatchString(link) {
		httpClient := http.Client{
			Timeout: time.Second * 10,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := httpClient.Get(link)
		if err != nil {
			return models.Coordinates{}, err
		}
		defer resp.Body.Close()

		location := resp.Header.Get("Location")

		return l.parseLink(location)
	}

	return l.parseLink(link)
}

func (l LocationLinkResolverImpl) parseLink(link string) (models.Coordinates, error) {
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
