package places

import (
	"context"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type Repository interface {
	SavePlace(ctx context.Context, place models.Place) error
	GetPlace(ctx context.Context, id string) (models.Place, error)
	GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, error)
	DeletePlace(ctx context.Context, id string) error
}

type Usecase interface {
	CreatePlace(ctx context.Context, request models.CreatePlaceRequest) error
	GetPlace(ctx context.Context, id string) (models.Place, error)
	GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, error)
	UpdatePlace(ctx context.Context, request models.UpdatePlaceRequest) error
	DeletePlace(ctx context.Context, id string) error
}
