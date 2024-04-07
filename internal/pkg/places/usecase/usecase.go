package usecase

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/googlesheets"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places"
)

type Usecase struct {
	repo         places.Repository
	linkResolver maps.LocationLinkResolver
	parser       googlesheets.SheetParser
	storageRepo  places.StorageRepository
}

func NewUsecase(repo places.Repository,
	linkResolver maps.LocationLinkResolver,
	parser googlesheets.SheetParser,
	storageRepo places.StorageRepository) Usecase {
	return Usecase{
		repo:         repo,
		linkResolver: linkResolver,
		parser:       parser,
		storageRepo:  storageRepo,
	}
}

func (u Usecase) SavePlace(ctx context.Context, place models.Place) models.Status {
	return u.repo.SavePlace(ctx, place)
}

func (u Usecase) GetPlace(ctx context.Context, id string) (models.Place, models.Status) {
	return u.repo.GetPlace(ctx, id)
}

func (u Usecase) GetPlaceByName(ctx context.Context, name string) (models.Place, models.Status) {
	return u.repo.GetPlaceByName(ctx, name)
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

func (u Usecase) GetCities(ctx context.Context) ([]models.City, models.Status) {
	return u.repo.GetCities(ctx)
}

func (u Usecase) GetCity(ctx context.Context, name string) (models.City, models.Status) {
	return u.repo.GetCity(ctx, name)
}

func (u Usecase) SaveCity(ctx context.Context, city models.City) models.Status {
	return u.repo.SaveCity(ctx, city)
}

func (u Usecase) GetPlacePhotoURLs(ctx context.Context, placeID string) ([]string, models.Status) {
	place, status := u.repo.GetPlace(ctx, placeID)
	if status.Code != models.OK {
		return nil, status
	}
	object := fmt.Sprintf("places/%s/%s", place.City, placeID)
	return u.storageRepo.GetPlacePhotoURLs(ctx, object)
}
