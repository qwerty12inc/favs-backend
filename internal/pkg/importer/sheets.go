package importer

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/mmcloughlin/geohash"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

func (i Importer) ImportCitiesFromSheet(ctx context.Context, sheetRange, city string) models.Status {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("recovering: %v", x)
		}
	}()

	places, status := i.parser.ParsePlaces(ctx, sheetRange)
	if status.Code != models.OK {
		return models.Status{models.InternalError, "failed to parse places"}
	}

	city = strings.ToLower(city)

	categoryLabels := make(map[string]map[string]struct{})
	for _, place := range places {
		if place.Category == "" {
			log.Printf("Skipping place without category: %s\n", place.Name)
			continue
		}
		if _, ok := categoryLabels[place.Category]; !ok {
			categoryLabels[place.Category] = make(map[string]struct{})
		}
		for _, label := range place.Labels {
			categoryLabels[place.Category][label] = struct{}{}
		}
	}

	cityInfo, err := i.linkResolver.GetCityInfo(ctx, city)
	if err != nil {
		log.Printf("models.Status while resolving city info: %v\n", err)
		return models.Status{models.InternalError, "failed to resolve city info"}
	}

	for category, labels := range categoryLabels {
		cityInfo.Categories = append(cityInfo.Categories, models.Category{
			Name: category,
			Labels: func() []string {
				var res []string
				for label := range labels {
					res = append(res, label)
				}
				return res
			}(),
		})
	}

	cityInfo.ImageURL = "places/" + city + "/city.jpg"
	status = i.service.SaveCity(ctx, cityInfo)
	if status.Code != models.OK {
		log.Printf("status while saving city: %v\n", status)
		return status
	}

	return models.Status{models.OK, "OK"}
}

func (i Importer) ImportPlacesFromSheet(ctx context.Context, sheetRange, city string) models.Status {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("recovering: %v", x)
		}
	}()

	places, status := i.parser.ParsePlaces(ctx, sheetRange)
	if status.Code != models.OK {
		return models.Status{models.InternalError, "failed to parse places"}
	}

	city = strings.ToLower(city)

	for _, place := range places {
		placeInfo, err := i.linkResolver.GetPlaceInfo(ctx, place.LocationURL, place.Name, city)
		if err != nil {
			log.Printf("models.Status while resolving coordinates: %v url: %s\n", err, place.LocationURL)
			return models.Status{models.InternalError, "failed to resolve coordinates"}
		}

		log.Printf("Getting place by name: %s\n", place.Name)
		oldPlace, status := i.service.GetPlaceByName(ctx, place.Name)
		if status.Code != models.OK {
			log.Printf("status while getting place by name: %v\n", status)
		}

		if oldPlace.ID != "" {
			placeInfo.ID = oldPlace.ID
		} else {
			placeInfo.ID = uuid.New().String()
		}

		placeInfo.GeoHash = geohash.Encode(placeInfo.Coordinates.Latitude, placeInfo.Coordinates.Longitude)
		placeInfo.Labels = place.Labels
		if placeInfo.City == "" {
			placeInfo.City = strings.ToLower(city)
		}
		if place.Instagram != "" {
			placeInfo.Instagram = place.Instagram
		}
		if place.Description != "" {
			placeInfo.Description = place.Description
		}
		if place.Website != "" {
			placeInfo.Website = place.Website
		}
		if place.LocationURL != "" {
			placeInfo.LocationURL = place.LocationURL
		}
		if place.Name != "" {
			placeInfo.Name = place.Name
		}

		log.Debug("Saving place: ", placeInfo.Name)
		status = i.service.SavePlace(ctx, *placeInfo)
		if status.Code != models.OK {
			log.Printf("error status while saving place: %v\n", status)
			return status
		}
		log.Printf("Place saved.\n")
	}

	return models.Status{models.OK, "OK"}
}
