package maps

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/v.rianov/favs-backend/internal/models"
	"googlemaps.github.io/maps"

	"net/http"
	"regexp"
	"strconv"
	"time"
)

type LocationLinkResolver interface {
	// ResolveLink resolves a location link and returns coordinates
	ResolveLink(link string) (models.Coordinates, error)
	GetPlaceInfo(ctx context.Context, link, name string) (*models.Place, error)
}

type LocationLinkResolverImpl struct {
	cl *maps.Client
}

var coordinatesRegexp = `@(-?\d+\.\d+),(-?\d+\.\d+)`
var googleMapsRegexp = `https:\/\/maps\.app\.goo\.gl\/[A-Za-z0-9]+`

func NewLocationLinkResolver(cl *maps.Client) LocationLinkResolverImpl {
	return LocationLinkResolverImpl{
		cl: cl,
	}
}

func (l LocationLinkResolverImpl) GetPlaceInfo(ctx context.Context, link, name string) (*models.Place, error) {
	log.Println("Resolving link: ", link)
	c, err := l.ResolveLink(link)
	if err != nil {
		log.Println("Error while resolving link: ", err)
		return nil, err
	}

	log.Println("Resolved coordinates: ", c)
	res, err := l.cl.TextSearch(ctx, &maps.TextSearchRequest{
		Query: name,
		Location: &maps.LatLng{
			Lat: c.Latitude,
			Lng: c.Longitude,
		},
		Radius: 2,
	})
	log.Println("Results: ", res.Results, " Error: ", err)
	if err != nil {
		log.Println("Error while searching for place: ", err)
		return nil, err
	}
	if len(res.Results) == 0 {
		log.Println("No places found for query: ", name, " coordinates: ", c)
		return nil, fmt.Errorf("no places found")
	}

	placeID := res.Results[0].PlaceID
	log.Println("Place ID: ", placeID)
	place, err := l.cl.PlaceDetails(ctx, &maps.PlaceDetailsRequest{
		PlaceID: placeID,
	})
	log.Println("Place details: ", place, " Error: ", err)
	if err != nil {
		log.Println("Error while getting place details: ", err)
		return nil, err
	}

	resPlace := &models.Place{
		Name:        place.Name,
		Description: place.FormattedAddress,
		LocationURL: place.URL,
		Coordinates: models.Coordinates{
			Latitude:  place.Geometry.Location.Lat,
			Longitude: place.Geometry.Location.Lng,
		},
		City:    "",
		Website: place.Website,
		Labels:  []string{},
		GeoHash: "",
		GoogleMapsInfo: &models.GoogleMapsInfo{
			PlaceID:          placeID,
			Rating:           place.Rating,
			Reservable:       place.Reservable,
			Delivery:         place.Delivery,
			FormattedAddress: place.FormattedAddress,
			LocationURL:      place.URL,
			Website:          place.Website,
			PhotoRefList: func() []string {
				var photoRefs []string
				for _, photo := range place.Photos {
					photoRefs = append(photoRefs, photo.PhotoReference)
				}
				return photoRefs
			}(),
		},
		Address: place.FormattedAddress,
	}

	if place.OpeningHours != nil {
		resPlace.GoogleMapsInfo.OpeningInfo = place.CurrentOpeningHours.WeekdayText
	}

	return resPlace, nil
}

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
