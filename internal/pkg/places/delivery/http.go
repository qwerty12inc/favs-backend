package delivery

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/usecase"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/utils"
)

type Handler struct {
	usecase usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

// GetPlace godoc
// @Summary Get place
// @Description Get place
// @Tags places
// @Accept json
// @Produce json
// @Param id path string true "Place ID"
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Success 200 {object} models.Place
// @Failure 404 "Place not found"
// @Failure 500 "Internal server error"
// @Router /places/{id} [get]
func (h Handler) GetPlace(c echo.Context) error {
	id := c.Param("id")
	place, status := h.usecase.GetPlace(c.Request().Context(), id)
	return utils.HandleResponse(c, status, place)
}

// TelegramGetPlace godoc
// @Summary Get place
// @Description Get place
// @Tags places
// @Accept json
// @Produce json
// @Param id path string true "Place ID"
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
//	@Param			X-Telegram-ID	header		string	true	"Telegram ID"
//
// @Success 200 {object} models.Place
// @Failure 404 "Place not found"
// @Failure 500 "Internal server error"
// @Router /tg/places/{id} [get]
func (h Handler) TelegramGetPlace(c echo.Context) error {
	id := c.Param("id")
	place, status := h.usecase.GetPlace(c.Request().Context(), id)
	return utils.HandleResponse(c, status, place)
}

// GetPlaces godoc
// @Summary Get places
// @Description Get places
// @Tags places
// @Accept json
// @Produce json
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Param latitude query float64 true "Latitude"
// @Param longitude query float64 true "Longitude"
// @Param latitudeDelta query float64 true "Latitude delta"
// @Param longitudeDelta query float64 true "Longitude delta"
// @Param labels query []string true "Labels"
// @Param category query string true "Category"
// @Param city query string true "City"
// @Success 200 {array} models.Place
// @Failure 400 "Invalid request"
// @Failure 404 "Place not found"
// @Failure 500 "Internal server error"
// @Router /places [get]
func (h Handler) GetPlaces(c echo.Context) error {
	request := models.GetPlacesRequest{}
	city := c.QueryParam("city")
	latitudeStr := c.QueryParam("latitude")
	longitudeStr := c.QueryParam("longitude")
	latitudeDeltaStr := c.QueryParam("latitudeDelta")
	longitudeDeltaStr := c.QueryParam("longitudeDelta")
	request.Category = c.QueryParam("category")
	request.Labels = c.QueryParams()["labels"]

	if city == "" {
		log.Debug("Getting places in box")
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			return err
		}
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			return err
		}
		latitudeDelta, err := strconv.ParseFloat(latitudeDeltaStr, 64)
		if err != nil {
			return err
		}
		longitudeDelta, err := strconv.ParseFloat(longitudeDeltaStr, 64)
		if err != nil {
			return err
		}
		request.Center = models.Coordinates{
			Latitude:  latitude,
			Longitude: longitude,
		}
		request.LatitudeDelta = latitudeDelta
		request.LongitudeDelta = longitudeDelta
	} else {
		log.Debug("Getting places in city")
		request.City = city
	}

	places, status := h.usecase.GetPlaces(c.Request().Context(), request)
	return utils.HandleResponse(c, status, places)
}

// TelegramGetPlaces godoc
// @Summary Get places
// @Description Get places
// @Tags places
// @Accept json
// @Produce json
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
//	@Param			X-Telegram-ID	header		string	true	"Telegram ID"
//
// @Param latitude query float64 true "Latitude"
// @Param longitude query float64 true "Longitude"
// @Param latitudeDelta query float64 true "Latitude delta"
// @Param longitudeDelta query float64 true "Longitude delta"
// @Param labels query []string true "Labels"
// @Param category query string true "Category"
// @Param city query string true "City"
// @Success 200 {array} models.Place
// @Failure 400 "Invalid request"
// @Failure 404 "Place not found"
// @Failure 500 "Internal server error"
// @Router /tg/places [get]
func (h Handler) TelegramGetPlaces(c echo.Context) error {
	request := models.GetPlacesRequest{}
	city := c.QueryParam("city")
	latitudeStr := c.QueryParam("latitude")
	longitudeStr := c.QueryParam("longitude")
	latitudeDeltaStr := c.QueryParam("latitudeDelta")
	longitudeDeltaStr := c.QueryParam("longitudeDelta")
	request.Category = c.QueryParam("category")
	request.Labels = c.QueryParams()["labels"]

	if city == "" {
		log.Debug("Getting places in box")
		latitude, err := strconv.ParseFloat(latitudeStr, 64)
		if err != nil {
			return err
		}
		longitude, err := strconv.ParseFloat(longitudeStr, 64)
		if err != nil {
			return err
		}
		latitudeDelta, err := strconv.ParseFloat(latitudeDeltaStr, 64)
		if err != nil {
			return err
		}
		longitudeDelta, err := strconv.ParseFloat(longitudeDeltaStr, 64)
		if err != nil {
			return err
		}
		request.Center = models.Coordinates{
			Latitude:  latitude,
			Longitude: longitude,
		}
		request.LatitudeDelta = latitudeDelta
		request.LongitudeDelta = longitudeDelta
	} else {
		log.Debug("Getting places in city")
		request.City = city
	}

	places, status := h.usecase.TelegramGetPlaces(c.Request().Context(), request)
	return utils.HandleResponse(c, status, places)
}

// TelegramGetCities godoc
// @Summary Get cities
// @Description Get cities
// @Tags cities
// @Accept json
// @Produce json
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
//	@Param			X-Telegram-ID	header		string	true	"Telegram ID"
//
// @Success 200 {array} models.City
// @Failure 404 "Cities not found"
// @Failure 500 "Internal server error"
// @Router /tg/cities [get]
func (h Handler) TelegramGetCities(c echo.Context) error {
	cities, status := h.usecase.GetCities(c.Request().Context())
	return utils.HandleResponse(c, status, cities)
}

// GetCities godoc
// @Summary Get cities
// @Description Get cities
// @Tags cities
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Success 200 {array} models.City
// @Failure 404 "Cities not found"
// @Failure 500 "Internal server error"
// @Router /cities [get]
func (h Handler) GetCities(c echo.Context) error {
	cities, status := h.usecase.GetCities(c.Request().Context())
	return utils.HandleResponse(c, status, cities)
}

// GetPlacePhotos godoc
// @Summary Get place photos
// @Description Get place photos
// @Tags places
// @Produce json
// @Param id path string true "Place ID"
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Success 200 {array} string
// @Failure 404 "Place not found"
// @Failure 500 "Internal server error"
// @Router /places/{id}/photos [get]
func (h Handler) GetPlacePhotos(c echo.Context) error {
	id := c.Param("id")
	photos, status := h.usecase.GetPlacePhotoURLs(c.Request().Context(), id)
	return utils.HandleResponse(c, status, photos)
}

func (h Handler) SaveUserPurchase(c echo.Context) error {
	purchaseStatus := c.QueryParam("status")
	if purchaseStatus != "success" {
		return c.JSON(400, "Invalid purchase status")
	}

	purchaseID := c.QueryParam("id")
	userEmail := c.QueryParam("user_email")
	amount := c.QueryParam("amount")
	price, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return c.JSON(400, "Invalid amount")
	}
	purchase := models.PurchaseObject{
		ID:    purchaseID,
		Price: int(price),
	}
	status := h.usecase.SaveUserPurchase(c.Request().Context(), userEmail, purchase)
	if status.Code != models.OK {
		return c.JSON(500, status)
	}
	return c.JSON(200, "OK")
}

// GeneratePaymentLink godoc
// @Summary Generate payment link
// @Description Generate payment link
// @Tags purchases
// @Produce json
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Param id query string true "Purchase ID"
// @Success 200 {text} string
// @Failure 400 "Invalid request"
// @Failure 500 "Internal server error"
// @Router /payments [get]
func (h Handler) GeneratePaymentLink(c echo.Context) error {
	purchaseID := c.QueryParam("id")
	userEmail := c.QueryParam("user_email")
	purchase := models.PurchaseObject{
		ID: purchaseID,
	}
	if userEmail == "" {
		user := c.Get("user").(*models.User)
		userEmail = user.Email
	}
	link, status := h.usecase.GeneratePaymentLink(c.Request().Context(), userEmail, purchase)
	return utils.HandleResponse(c, status, link)
}

// SaveReport godoc
// @Summary Save report
// @Description Save report
// @Tags reports
// @Produce json
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Param id path string true "Place ID"
// @Success 200 {text} string
// @Failure 400 "Invalid request"
// @Failure 500 "Internal server error"
// @Router /places/{id}/reports [post]
func (h Handler) SaveReport(c echo.Context) error {
	report := models.Report{}
	err := c.Bind(&report)
	if err != nil {
		return c.JSON(400, "Invalid request")
	}
	report.PlaceID = c.Param("id")
	status := h.usecase.SaveReport(c.Request().Context(), report)
	return utils.HandleResponse(c, status, "OK")
}
