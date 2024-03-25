package delivery

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
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
// @Param city query string true "City"
// @Success 200 {array} models.Place
// @Failure 500 "Internal server error"
// @Router /places [get]
func (h Handler) GetPlaces(c echo.Context) error {
	request := models.GetPlacesRequest{}
	city := c.QueryParam("city")
	latitudeStr := c.QueryParam("latitude")
	longitudeStr := c.QueryParam("longitude")
	latitudeDeltaStr := c.QueryParam("latitudeDelta")
	longitudeDeltaStr := c.QueryParam("longitudeDelta")
	request.Labels = c.QueryParams()["labels"]

	if city == "" {
		log.Println("Getting places in box")
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
		log.Println("Getting places in city")
		request.City = city
	}

	places, status := h.usecase.GetPlaces(c.Request().Context(), request)
	return utils.HandleResponse(c, status, places)
}

// GetCities godoc
// @Summary Get cities
// @Description Get cities
// @Tags cities
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Success 200 {array} string
// @Failure 500 "Internal server error"
// @Router /cities [get]
func (h Handler) GetCities(c echo.Context) error {
	cities, status := h.usecase.GetCities(c.Request().Context())
	return utils.HandleResponse(c, status, cities)
}

// GetLabels godoc
// @Summary Get labels
// @Description Get labels
// @Tags labels
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Success 200 {array} string
// @Failure 500 "Internal server error"
// @Router /labels [get]
func (h Handler) GetLabels(c echo.Context) error {
	labels, status := h.usecase.GetLabels(c.Request().Context())
	return utils.HandleResponse(c, status, labels)
}
