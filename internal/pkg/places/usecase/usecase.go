package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/mmcloughlin/geohash"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/googlesheets"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places"
)

type Usecase struct {
	repo         places.Repository
	linkResolver maps.LocationLinkResolver
	parser       googlesheets.SheetParser
}

func NewUsecase(repo places.Repository,
	linkResolver maps.LocationLinkResolver,
	parser googlesheets.SheetParser) Usecase {
	return Usecase{
		repo:         repo,
		linkResolver: linkResolver,
		parser:       parser,
	}
}

func (u Usecase) CreatePlace(ctx context.Context, request models.CreatePlaceRequest) models.Status {
	coordinates, err := u.linkResolver.ResolveLink(request.LocationURL)
	if err != nil {
		return models.Status{models.BadRequest, fmt.Sprintf("failed to resolve coordinates: %v", err)}
	}

	place := models.Place{
		Name:        request.Name,
		Description: request.Description,
		LocationURL: request.LocationURL,
		Coordinates: coordinates,
		City:        request.City,
		Website:     request.Website,
		Labels:      request.Labels,
		GeoHash:     geohash.Encode(coordinates.Latitude, coordinates.Longitude),
	}
	return u.repo.SavePlace(ctx, place)
}

func (u Usecase) ImportPlacesFromSheet(ctx context.Context, sheetRange string,
	city string, force bool) models.Status {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("recovering: %v", x)
		}
	}()
	places, status := u.parser.ParsePlaces(ctx, sheetRange)
	if status.Code != models.OK {
		return models.Status{models.InternalError, "failed to parse places"}
	}

	for _, place := range places {
		log.Printf("Getting place by name: %s\n", place.Name)
		oldPlace, status := u.repo.GetPlaceByName(ctx, place.Name)
		if status.Code == models.OK && oldPlace.Name == place.Name && !force {
			log.Printf("Place with name %s already exists, skipping.\n", place.Name)
			continue
		}

		placeInfo, err := u.linkResolver.GetPlaceInfo(ctx, place.LocationURL, place.Name)
		if err != nil {
			log.Printf("models.Status while resolving coordinates: %v url: %s\n", err, place.LocationURL)
			return models.Status{models.InternalError, "failed to resolve coordinates"}
		}

		placeInfo.ID = uuid.New().String()
		placeInfo.GeoHash = geohash.Encode(placeInfo.Coordinates.Latitude, placeInfo.Coordinates.Longitude)
		placeInfo.Labels = place.Labels
		if placeInfo.City == "" {
			placeInfo.City = strings.ToLower(city)
		}
		if placeInfo.Instagram == "" {
			placeInfo.Instagram = place.Instagram
		}
		if placeInfo.Description == "" {
			placeInfo.Description = place.Description
		}
		if placeInfo.Website == "" {
			placeInfo.Website = place.Website
		}
		if placeInfo.LocationURL == "" {
			placeInfo.LocationURL = place.LocationURL
		}
		if placeInfo.Name == "" {
			placeInfo.Name = place.Name
		}

		log.Println("Saving place: ", placeInfo.Name)
		status = u.repo.SavePlace(ctx, *placeInfo)
		if status.Code != models.OK {
			log.Printf("models.Status while saving place: %v\n", status)
			return status
		}
		log.Printf("Place saved.\n")
	}
	return models.Status{models.OK, "OK"}
}

func (u Usecase) GetPlace(ctx context.Context, id string) (models.Place, models.Status) {
	return u.repo.GetPlace(ctx, id)
}

func (u Usecase) GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status) {
	request.City = strings.ToLower(request.City)
	places, status := u.repo.GetPlaces(ctx, request)
	if status.Code != models.OK {
		return nil, status
	}

	if len(places) == 0 {
		return nil, models.Status{models.NotFound, "places not found"}
	}
	return places, models.Status{models.OK, "OK"}
}

func (u Usecase) UpdatePlace(ctx context.Context, request models.UpdatePlaceRequest) models.Status {
	place, status := u.repo.GetPlace(ctx, request.ID)
	if status.Code != models.OK {
		return status
	}
	return u.repo.SavePlace(ctx, place)
}

func (u Usecase) DeletePlace(ctx context.Context, id string) models.Status {
	return u.repo.DeletePlace(ctx, id)
}

func (u Usecase) GetCities(ctx context.Context) ([]string, models.Status) {
	return u.repo.GetCities(ctx)
}
