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
	GetCities(ctx context.Context) ([]models.City, models.Status)
	SaveCity(ctx context.Context, city models.City) models.Status
	GetCity(ctx context.Context, name string) (models.City, models.Status)
	GetUserPurchases(ctx context.Context, userEmail string) (models.UserPurchases, models.Status)
	SaveUserPurchase(ctx context.Context, userEmail string, purchase models.PurchaseObject) models.Status
}

type Usecase interface {
	SavePlace(ctx context.Context, place models.Place) models.Status
	GetPlace(ctx context.Context, id string) (models.Place, models.Status)
	GetPlaceByName(ctx context.Context, name string) (models.Place, models.Status)
	GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status)
	GetCities(ctx context.Context) ([]models.City, models.Status)
	SaveCity(ctx context.Context, city models.City) models.Status
	GetCity(ctx context.Context, name string) (models.City, models.Status)
	SaveUserPurchase(ctx context.Context, userEmail string, purchase models.PurchaseObject) models.Status
	GeneratePaymentLink(ctx context.Context, userEmail string, purchase models.PurchaseObject) (string, models.Status)
	TelegramGetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status)
}

type StorageRepository interface {
	GetPlacePhotoURLs(ctx context.Context, object string) ([]string, models.Status)
	GenerateSignedURL(ctx context.Context, object string) (string, models.Status)
}
