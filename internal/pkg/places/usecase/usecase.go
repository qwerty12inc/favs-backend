package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/googlesheets"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places"
	"log"
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

func (u Usecase) CreatePlace(ctx context.Context, request models.CreatePlaceRequest) error {
	coordinates, err := u.linkResolver.ResolveLink(request.LocationURL)
	if err != nil {
		return err
	}

	place := models.Place{
		Name:        request.Name,
		Description: request.Description,
		LocationURL: request.LocationURL,
		Coordinates: coordinates,
		City:        request.City,
		Website:     request.Website,
		Labels:      request.Labels,
	}
	return u.repo.SavePlace(ctx, place)
}

func (u Usecase) ImportPlacesFromSheet(ctx context.Context, sheetRange string,
	city string, force bool) error {
	places, status := u.parser.ParsePlaces(ctx, sheetRange)
	if status.Code != models.OK {
		return fmt.Errorf("failed to parse places from sheet: %s", status.Message)
	}

	for _, place := range places {
		coordinates, err := u.linkResolver.ResolveLink(place.LocationURL)
		if err != nil {
			log.Printf("Error while resolving coordinates: %v url: %s\n", err, place.LocationURL)
		}

		placeModel := models.Place{
			ID:          uuid.New().String(),
			Name:        place.Name,
			Description: place.Description,
			LocationURL: place.LocationURL,
			Coordinates: coordinates,
			City:        city,
			Website:     place.Website,
			Instagram:   place.Instagram,
			Labels:      place.Labels,
		}

		oldPlace, err := u.repo.GetPlaceByName(ctx, placeModel.Name)
		if err == nil && oldPlace.Name == placeModel.Name && !force {
			log.Printf("Place with name %s already exists, skipping.\n", placeModel.Name)
			continue
		}

		log.Printf("Saving place: %v\n", placeModel)

		err = u.repo.SavePlace(ctx, placeModel)
		if err != nil {
			log.Printf("Error while saving place: %v\n", err)
			return err
		}
		log.Printf("Place saved.\n")
	}
	return nil
}

func (u Usecase) GetPlace(ctx context.Context, id string) (models.Place, error) {
	return u.repo.GetPlace(ctx, id)
}

func (u Usecase) GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, error) {
	return u.repo.GetPlaces(ctx, request)
}

func (u Usecase) UpdatePlace(ctx context.Context, request models.UpdatePlaceRequest) error {
	place, err := u.repo.GetPlace(ctx, request.ID)
	if err != nil {
		return err
	}
	return u.repo.SavePlace(ctx, place)
}

func (u Usecase) DeletePlace(ctx context.Context, id string) error {
	return u.repo.DeletePlace(ctx, id)
}
