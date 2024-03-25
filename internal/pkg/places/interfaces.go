package places

import (
	"context"

	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type Repository interface {
	SavePlace(ctx context.Context, place models.Place) models.Status
	GetPlace(ctx context.Context, id string) (models.Place, models.Status)
	GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status)
	DeletePlace(ctx context.Context, id string) models.Status
	GetPlaceByName(ctx context.Context, name string) (models.Place, models.Status)
	GetCities(ctx context.Context) ([]string, models.Status)
	GetLabels(ctx context.Context) ([]string, models.Status)
}

type Usecase interface {
	CreatePlace(ctx context.Context, request models.CreatePlaceRequest) models.Status
	GetPlace(ctx context.Context, id string) (models.Place, models.Status)
	GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status)
	UpdatePlace(ctx context.Context, request models.UpdatePlaceRequest) models.Status
	DeletePlace(ctx context.Context, id string) models.Status
	ImportPlacesFromSheet(ctx context.Context, sheetRange string, city string, force bool) models.Status
	GetCities(ctx context.Context) ([]string, models.Status)
	GetLabels(ctx context.Context) ([]string, models.Status)
}
