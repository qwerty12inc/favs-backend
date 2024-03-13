package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mmcloughlin/geohash"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/googlesheets"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places"
	"log"
	"strings"
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
	places, status := u.parser.ParsePlaces(ctx, sheetRange)
	if status.Code != models.OK {
		return models.Status{models.InternalError, "failed to parse places"}
	}

	for _, place := range places {
		coordinates, err := u.linkResolver.ResolveLink(place.LocationURL)
		if err != nil {
			log.Printf("models.Status while resolving coordinates: %v url: %s\n", err, place.LocationURL)
		}

		log.Printf("Coordinates resolved: %v\n", coordinates)

		placeModel := models.Place{
			ID:          uuid.New().String(),
			Name:        place.Name,
			Description: place.Description,
			LocationURL: place.LocationURL,
			Coordinates: coordinates,
			City:        strings.ToLower(city),
			Website:     place.Website,
			Instagram:   place.Instagram,
			Labels:      place.Labels,
			GeoHash:     geohash.Encode(coordinates.Latitude, coordinates.Longitude),
		}

		oldPlace, status := u.repo.GetPlaceByName(ctx, placeModel.Name)
		if status.Code == models.OK && oldPlace.Name == placeModel.Name && !force {
			log.Printf("Place with name %s already exists, skipping.\n", placeModel.Name)
			continue
		}

		log.Printf("Saving place: %v\n", placeModel)

		status = u.repo.SavePlace(ctx, placeModel)
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
