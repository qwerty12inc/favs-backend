package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"google.golang.org/api/iterator"
)

type Repository struct {
	cl *firestore.Client
}

func NewRepository(cl *firestore.Client) Repository {
	return Repository{
		cl: cl,
	}
}

func (r Repository) SavePlace(ctx context.Context, place models.Place) error {
	_, err := r.cl.Collection("places").Doc(place.ID).Create(ctx, place)
	return err
}

func (r Repository) GetPlace(ctx context.Context, id string) (models.Place, error) {
	doc, err := r.cl.Collection("places").Doc(id).Get(ctx)
	if err != nil {
		return models.Place{}, err
	}
	var place models.Place
	err = doc.DataTo(&place)
	if err != nil {
		return models.Place{}, err
	}
	return place, nil
}

func (r Repository) GetPlaceByName(ctx context.Context, name string) (models.Place, error) {
	iter := r.cl.Collection("places").Where("name", "==", name).Documents(ctx)
	doc, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return models.Place{}, errors.New("place not found")
	}
	if err != nil {
		return models.Place{}, err
	}
	var place models.Place
	err = doc.DataTo(&place)
	if err != nil {
		return models.Place{}, err
	}
	return place, nil
}

func (r Repository) GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, error) {
	iter := r.cl.Collection("places").
		Where("coordinates/latitude", ">", request.Center.Latitude-request.LatitudeDelta).
		Where("coordinates/latitude", "<", request.Center.Latitude+request.LatitudeDelta).
		Where("coordinates/longitude", ">", request.Center.Longitude-request.LongitudeDelta).
		Where("coordinates/longitude", "<", request.Center.Longitude+request.LongitudeDelta).
		Documents(ctx)

	var places []models.Place
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		var place models.Place
		err = doc.DataTo(&place)
		if err != nil {
			return nil, err
		}
		places = append(places, place)
	}
	return places, nil
}

func (r Repository) DeletePlace(ctx context.Context, id string) error {
	_, err := r.cl.Collection("places").Doc(id).Delete(ctx)
	return err
}
