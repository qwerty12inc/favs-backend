package usecase

import (
	"context"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places"
)

type Usecase struct {
	repo         places.Repository
	linkResolver maps.LocationLinkResolver
}

func NewUsecase(repo places.Repository,
	linkResolver maps.LocationLinkResolver) Usecase {
	return Usecase{
		repo:         repo,
		linkResolver: linkResolver,
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
		OpenAt:      request.OpenAt,
		ClosedAt:    request.ClosedAt,
		City:        request.City,
		Address:     request.Address,
		Phone:       request.Phone,
		Type:        request.Type,
		Website:     request.Website,
		Labels:      request.Labels,
	}
	return u.repo.SavePlace(ctx, place)
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
