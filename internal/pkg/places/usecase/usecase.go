package usecase

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
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
	place, status := u.repo.GetPlace(ctx, id)
	if status.Code != models.OK {
		return models.Place{}, status
	}

	// transform all photo references to signed urls
	if place.GoogleMapsInfo != nil && len(place.GoogleMapsInfo.PhotoRefList) != 0 {
		for i, ref := range place.GoogleMapsInfo.PhotoRefList {
			place.GoogleMapsInfo.PhotoRefList[i], status = u.storageRepo.GenerateSignedURL(ctx, ref)
			if status.Code != models.OK {
				log.Error("Failed to generate signed URL ", status)
			}
		}
	}

	place.IsOpen = place.IsOpenNow()

	return place, models.Status{Code: models.OK, Message: "OK"}
}

func (u Usecase) GetPlaceByName(ctx context.Context, name string) (models.Place, models.Status) {
	return u.repo.GetPlaceByName(ctx, name)
}

func (u Usecase) getPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status) {
	places, status := u.repo.GetPlaces(ctx, request)
	if status.Code != models.OK {
		log.Error("Failed to get places ", status)
		return nil, status
	}

	start := time.Now()
	wg := sync.WaitGroup{}
	// transform all photo references to signed urls
	for i := range places {
		if places[i].GoogleMapsInfo != nil && len(places[i].GoogleMapsInfo.PhotoRefList) != 0 {
			for j, ref := range places[i].GoogleMapsInfo.PhotoRefList {
				wg.Add(1)
				go func(i, j int, ref string) {
					defer wg.Done()
					places[i].GoogleMapsInfo.PhotoRefList[j], status = u.storageRepo.GenerateSignedURL(ctx, ref)
					if status.Code != models.OK {
						log.Error("Failed to generate signed URL ", status)
					}
				}(i, j, ref)
			}
		}
	}
	wg.Wait()
	log.Info("Time to generate signed URLs: ", time.Since(start))

	for i := range places {
		if places[i].ImagePreview == "" &&
			places[i].GoogleMapsInfo != nil &&
			len(places[i].GoogleMapsInfo.PhotoRefList) != 0 {
			places[i].ImagePreview = places[i].GoogleMapsInfo.PhotoRefList[0]
		}
		// define whether the place is open now
		if places[i].GoogleMapsInfo != nil && places[i].GoogleMapsInfo.OpeningInfo != nil {
			places[i].IsOpen = places[i].IsOpenNow()
		}
	}

	if len(places) == 0 {
		log.Error("Places not found")
		return nil, models.Status{models.NotFound, "places not found"}
	}
	return places, models.Status{models.OK, "OK"}
}

func (u Usecase) GetPlaces(ctx context.Context, request models.GetPlacesRequest) ([]models.Place, models.Status) {
	request.City = strings.ToLower(request.City)
	city, status := u.repo.GetCity(ctx, request.City)
	if status.Code != models.OK {
		log.Error("Failed to get city ", status)
		return nil, status
	}

	needsPurchase := false
	stripeProductID := ""
	for _, category := range city.Categories {
		if category.Name == request.Category && category.NeedsPurchase {
			needsPurchase = true
			stripeProductID = category.StripeProductID
			break
		}
	}

	if needsPurchase {
		user := ctx.Value("user").(models.User)
		userPurchases, status := u.repo.GetUserPurchases(ctx, user.Email)
		if !userPurchases.HasPurchase(stripeProductID) || status.Code != models.OK {
			link, status := u.GeneratePaymentLink(ctx, user.Email, models.PurchaseObject{
				ID: stripeProductID,
			})
			if status.Code != models.OK {
				log.Error("Failed to generate payment link ", status)
				return nil, status
			}
			return nil, models.Status{models.Forbidden, "You need to purchase this category. Payment link: " + link}
		}
	}

	return u.getPlaces(ctx, request)
}

func (u Usecase) TelegramGetPlaces(ctx context.Context,
	request models.GetPlacesRequest) ([]models.Place, models.Status) {
	res, status := u.getPlaces(ctx, request)
	if status.Code != models.OK {
		return nil, status
	}
	if len(res) <= 5 {
		return res, status
	}
	return res[:5], status
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
	pr, err := u.stripeConnector.GetProductByID(purchase.ID)
	if err != nil {
		log.Error("Failed to get product ", err)
		return "", models.Status{models.InternalError, err.Error()}
	}

	link, err := paymentlink.New(&stripe.PaymentLinkParams{
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				Quantity: stripe.Int64(1),
				Price:    stripe.String(pr.PriceID),
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
		log.Error("Failed to create payment link ", err)
		return "", models.Status{models.InternalError, err.Error()}
	}
	return link.URL, models.Status{models.OK, "OK"}
}

func (u Usecase) SaveReport(ctx context.Context, report models.Report) models.Status {
	report.ReportedAt = time.Now().Unix()
	report.ID = fmt.Sprintf("%s-ts-%d", report.ReportedBy, report.ReportedAt)
	return u.repo.SaveReport(ctx, report)
}

func (u Usecase) GetReports(ctx context.Context) ([]models.Report, models.Status) {
	return u.repo.GetReports(ctx)
}

func (u Usecase) AddUserPlace(ctx context.Context, request models.AddPlaceRequest) models.Status {
	request.AddedAt = time.Now().Unix()
	request.ID = fmt.Sprintf("%s-ts-%d", request.AddedBy, request.AddedAt)
	return u.repo.AddUserPlace(ctx, request)
}
