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
	_, err := r.cl.Collection("places").Doc(place.ID).Set(ctx, place)
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

func (r Repository) GetCities(ctx context.Context) ([]models.City, models.Status) {
	iter := r.cl.Collection("cities").Documents(ctx)

	var cities []models.City
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, models.Status{models.InternalError, err.Error()}
		}
		var city models.City
		err = doc.DataTo(&city)
		if err != nil {
			return nil, models.Status{models.InternalError, err.Error()}
		}
		cities = append(cities, city)
	}
	return cities, models.Status{models.OK, "OK"}
}

func (r Repository) SaveCity(ctx context.Context, city models.City) models.Status {
	_, err := r.cl.Collection("cities").Doc(city.Name).Set(ctx, city)
	if err != nil {
		return models.Status{models.InternalError, err.Error()}
	}
	return models.Status{models.OK, "OK"}
}

func (r Repository) GetCity(ctx context.Context, name string) (models.City, models.Status) {
	doc, err := r.cl.Collection("cities").Doc(name).Get(ctx)
	if err != nil {
		return models.City{}, models.Status{models.InternalError, err.Error()}
	}
	var city models.City
	err = doc.DataTo(&city)
	if err != nil {
		return models.City{}, models.Status{models.InternalError, err.Error()}
	}
	return city, models.Status{models.OK, "OK"}
}

func (r Repository) GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status) {
	var iter *firestore.DocumentIterator
	var query firestore.Query
	if request.City == "" {
		box := geohash.Box{
			MinLat: request.Center.Latitude - request.LatitudeDelta,
			MaxLat: request.Center.Latitude + request.LatitudeDelta,
			MinLng: request.Center.Longitude - request.LongitudeDelta,
			MaxLng: request.Center.Longitude + request.LongitudeDelta,
		}
		log.Println("Getting places in box: ", box)
		log.Println("Min: ", geohash.Encode(box.MinLat, box.MinLng))
		log.Println("Max: ", geohash.Encode(box.MaxLat, box.MaxLng))
		query = r.cl.Collection("places").OrderBy("geohash", firestore.Asc).
			StartAt(geohash.Encode(box.MinLat, box.MinLng)).
			EndAt(geohash.Encode(box.MaxLat, box.MaxLng))
	} else {
		log.Println("Getting places in city: ", request.City)
		query = r.cl.Collection("places").
			Where("city", "==", request.City)
	}

	log.Println("Length of labels: ", len(request.Labels))

	if request.Category != "" {
		log.Println("Filtering by category: ", request.Category)
		query = query.Where("category", "==", request.Category)
	}

	if len(request.Labels) > 0 {
		log.Println("Filtering by labels: ", request.Labels)
		query = query.Where("labels", "array-contains-any", request.Labels)
	}

	iter = query.Documents(ctx)

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

func (r Repository) GetLabels(ctx context.Context, city string) ([]string, models.Status) {
	query := r.cl.Collection("places").Select("labels")
	if city != "" {
		query = query.Where("city", "==", city)
	}
	iter := query.Documents(ctx)

	labels := make(map[string]bool)
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
		for _, label := range place.Labels {
			labels[label] = true
		}
	}
	filteredLabels := make([]string, 0, len(labels))
	for label := range labels {
		filteredLabels = append(filteredLabels, label)
	}
	return filteredLabels, models.Status{models.OK, "OK"}
}
