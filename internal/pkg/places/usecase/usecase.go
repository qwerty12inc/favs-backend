package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mmcloughlin/geohash"
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
	city, category string, force bool) models.Status {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("recovering: %v", x)
		}
	}()

	places, status := u.parser.ParsePlaces(ctx, sheetRange)
	if status.Code != models.OK {
		return models.Status{models.InternalError, "failed to parse places"}
	}

	city = strings.ToLower(city)

	categoryLabels := make(map[string]map[string]struct{})
	for _, place := range places {
		if place.Category == "" {
			log.Printf("Skipping place without category: %s\n", place.Name)
			continue
		}
		if _, ok := categoryLabels[place.Category]; !ok {
			categoryLabels[place.Category] = make(map[string]struct{})
		}
		for _, label := range place.Labels {
			categoryLabels[place.Category][label] = struct{}{}
		}
	}

	//for _, place := range places {
	//	log.Printf("Getting place by name: %s\n", place.Name)
	//	oldPlace, status := u.repo.GetPlaceByName(ctx, place.Name)
	//
	//	placeInfo, err := u.linkResolver.GetPlaceInfo(ctx, place.LocationURL, place.Name, city)
	//	if err != nil {
	//		log.Printf("models.Status while resolving coordinates: %v url: %s\n", err, place.LocationURL)
	//		return models.Status{models.InternalError, "failed to resolve coordinates"}
	//	}
	//
	//	if oldPlace.ID != "" {
	//		placeInfo.ID = oldPlace.ID
	//	} else {
	//		placeInfo.ID = uuid.New().String()
	//	}
	//
	//	placeInfo.GeoHash = geohash.Encode(placeInfo.Coordinates.Latitude, placeInfo.Coordinates.Longitude)
	//	placeInfo.Labels = place.Labels
	//	if placeInfo.City == "" {
	//		placeInfo.City = strings.ToLower(city)
	//	}
	//	if placeInfo.Instagram == "" {
	//		placeInfo.Instagram = place.Instagram
	//	}
	//	if placeInfo.Description == "" {
	//		placeInfo.Description = place.Description
	//	}
	//	if placeInfo.Website == "" {
	//		placeInfo.Website = place.Website
	//	}
	//	if placeInfo.LocationURL == "" {
	//		placeInfo.LocationURL = place.LocationURL
	//	}
	//	if placeInfo.Name == "" {
	//		placeInfo.Name = place.Name
	//	}
	//	if placeInfo.Category == "" {
	//		placeInfo.Category = category
	//	}
	//
	//	log.Println("Saving place: ", placeInfo.Name)
	//	status = u.repo.SavePlace(ctx, *placeInfo)
	//	if status.Code != models.OK {
	//		log.Printf("models.Status while saving place: %v\n", status)
	//		return status
	//	}
	//	log.Printf("Place saved.\n")
	//}

	cityInfo, err := u.linkResolver.GetCityInfo(ctx, city)
	if err != nil {
		log.Printf("models.Status while resolving city info: %v\n", err)
		return models.Status{models.InternalError, "failed to resolve city info"}
	}

	for category, labels := range categoryLabels {
		cityInfo.Categories = append(cityInfo.Categories, models.Category{
			Name: category,
			Labels: func() []string {
				var res []string
				for label := range labels {
					res = append(res, label)
				}
				return res
			}(),
		})
	}

	cityInfo.ImageURL = "places/" + city + "/city.jpg"
	status = u.repo.SaveCity(ctx, cityInfo)
	if status.Code != models.OK {
		log.Printf("status while saving city: %v\n", status)
		return status
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

func (u Usecase) GetCities(ctx context.Context) ([]models.City, models.Status) {
	return u.repo.GetCities(ctx)
}

func (u Usecase) SaveCity(ctx context.Context, city models.City) models.Status {
	return u.repo.SaveCity(ctx, city)
}

func (u Usecase) GetCity(ctx context.Context, name string) (models.City, models.Status) {
	return u.repo.GetCity(ctx, name)
}

func (u Usecase) GetPlacePhotoURLs(ctx context.Context, placeID string) ([]string, models.Status) {
	place, status := u.repo.GetPlace(ctx, placeID)
	if status.Code != models.OK {
		return nil, status
	}
	object := fmt.Sprintf("places/%s/%s", place.City, placeID)
	return u.storageRepo.GetPlacePhotoURLs(ctx, object)
}
