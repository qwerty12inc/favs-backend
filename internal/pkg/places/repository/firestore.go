package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/mmcloughlin/geohash"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"google.golang.org/api/iterator"
	"log"
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
	log.Println("Saving place: ", place)
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
	box := geohash.Box{
		MinLat: request.Center.Latitude - request.LatitudeDelta,
		MaxLat: request.Center.Latitude + request.LatitudeDelta,
		MinLng: request.Center.Longitude - request.LongitudeDelta,
		MaxLng: request.Center.Longitude + request.LongitudeDelta,
	}
	var iter *firestore.DocumentIterator
	if request.City == "" {
		log.Println("Getting places in box: ", box)
		log.Println("Min: ", geohash.Encode(box.MinLat, box.MinLng))
		log.Println("Max: ", geohash.Encode(box.MaxLat, box.MaxLng))
		iter = r.cl.Collection("places").OrderBy("geohash", firestore.Asc).
			StartAt(geohash.Encode(box.MinLat, box.MinLng)).
			EndAt(geohash.Encode(box.MaxLat, box.MaxLng)).Documents(ctx)
	} else {
		iter = r.cl.Collection("places").
			Where("city", "==", request.City).
			Documents(ctx)
	}

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
