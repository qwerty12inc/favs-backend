package repository

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/mmcloughlin/geohash"
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

func (r Repository) SavePlace(ctx context.Context, place models.Place) models.Status {
	log.Println("Saving place: ", place)
	_, err := r.cl.Collection("places").Doc(place.ID).Create(ctx, place)
	if err != nil {
		return models.Status{models.InternalError, err.Error()}
	}
	return models.Status{models.OK, "OK"}
}

func (r Repository) GetPlace(ctx context.Context, id string) (models.Place, models.Status) {
	doc, err := r.cl.Collection("places").Doc(id).Get(ctx)
	if err != nil {
		return models.Place{}, models.Status{models.InternalError, err.Error()}
	}
	var place models.Place
	err = doc.DataTo(&place)
	if err != nil {
		return models.Place{}, models.Status{models.InternalError, err.Error()}
	}
	return place, models.Status{models.OK, "OK"}
}

func (r Repository) GetPlaceByName(ctx context.Context, name string) (models.Place, models.Status) {
	iter := r.cl.Collection("places").Where("name", "==", name).Documents(ctx)
	doc, err := iter.Next()
	if errors.Is(err, iterator.Done) {
		return models.Place{}, models.Status{models.NotFound, "Place not found"}
	}
	if err != nil {
		return models.Place{}, models.Status{models.InternalError, err.Error()}
	}
	var place models.Place
	err = doc.DataTo(&place)
	if err != nil {
		return models.Place{}, models.Status{models.InternalError, err.Error()}
	}
	return place, models.Status{models.OK, "OK"}
}

func (r Repository) GetCities(ctx context.Context) ([]string, models.Status) {
	iter := r.cl.Collection("places").Select("city").Documents(ctx)
	var cities map[string]bool
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, models.Status{models.InternalError, err.Error()}
		}
		var place models.Place
		err = doc.DataTo(&place)
		if err != nil {
			return nil, models.Status{models.InternalError, err.Error()}
		}
		cities[place.City] = true
	}
	filteredCities := make([]string, 0, len(cities))
	for city := range cities {
		filteredCities = append(filteredCities, city)
	}
	return filteredCities, models.Status{models.OK, "OK"}
}

func (r Repository) GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status) {
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
			return nil, models.Status{models.InternalError, err.Error()}
		}
		var place models.Place
		err = doc.DataTo(&place)
		if err != nil {
			return nil, models.Status{models.InternalError, err.Error()}
		}
		places = append(places, place)
	}
	return places, models.Status{models.OK, "OK"}
}

func (r Repository) DeletePlace(ctx context.Context, id string) models.Status {
	_, err := r.cl.Collection("places").Doc(id).Delete(ctx)
	if err != nil {
		return models.Status{models.InternalError, err.Error()}
	}
	return models.Status{models.OK, "OK"}
}
