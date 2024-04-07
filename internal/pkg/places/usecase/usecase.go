package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentlink"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/googlesheets"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places"
	stripe2 "gitlab.com/v.rianov/favs-backend/internal/pkg/stripe"
)

type Usecase struct {
	repo            places.Repository
	linkResolver    maps.LocationLinkResolver
	parser          googlesheets.SheetParser
	storageRepo     places.StorageRepository
	stripeConnector stripe2.StripeConnector
}

func NewUsecase(repo places.Repository,
	linkResolver maps.LocationLinkResolver,
	parser googlesheets.SheetParser,
	storageRepo places.StorageRepository,
	stripeConnector stripe2.StripeConnector) Usecase {
	return Usecase{
		repo:            repo,
		linkResolver:    linkResolver,
		parser:          parser,
		storageRepo:     storageRepo,
		stripeConnector: stripeConnector,
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
	city, status := u.repo.GetCity(ctx, request.City)
	if status.Code != models.OK {
		return nil, status
	}

	needsPurchase := false
	for _, category := range city.Categories {
		if category.Name == request.Category && category.NeedsPurchase {
			needsPurchase = true
			break
		}
	}

	if needsPurchase {
		user := ctx.Value("user").(models.User)
		userPurchases, status := u.repo.GetUserPurchases(ctx, user.Email)
		if status.Code != models.OK {
			return nil, status
		}
		if !userPurchases.HasPurchase(request.Category) {
			return nil, models.Status{models.Forbidden, "You need to purchase this category"}
		}
	}

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

func (u Usecase) SaveUserPurchase(ctx context.Context, userEmail string, purchase models.PurchaseObject) models.Status {
	return u.repo.SaveUserPurchase(ctx, userEmail, purchase)
}

func (u Usecase) GeneratePaymentLink(ctx context.Context, userEmail string, purchase models.PurchaseObject) (string, models.Status) {
	pr, err := u.stripeConnector.GetProductByName(purchase.ID)
	if err != nil {
		return "", models.Status{models.InternalError, err.Error()}
	}

	link, err := paymentlink.New(&stripe.PaymentLinkParams{
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				ID:       stripe.String(purchase.ID),
				Quantity: stripe.Int64(pr.Price),
			},
		},
		Currency: stripe.String("usd"),
		AfterCompletion: &stripe.PaymentLinkAfterCompletionParams{
			// TODO: add a code to the URL to check if the purchase was successful
			Redirect: &stripe.PaymentLinkAfterCompletionRedirectParams{
				URL: stripe.String("https://favs.site/api/v1/purchases?status=success&id=" +
					purchase.ID + "&user_email=" + userEmail + "&amount=" + fmt.Sprint(purchase.Price)),
			},
			Type: stripe.String("redirect"),
		},
	})
	if err != nil {
		return "", models.Status{models.InternalError, err.Error()}
	}
	return link.URL, models.Status{models.OK, "OK"}
}
