package delivery

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/usecase"
	"strconv"
)

type Handler struct {
	usecase usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

// CreatePlace godoc
// @Summary Create place
// @Description Create place
// @Tags places
// @Accept json
// @Produce json
// @Param place body models.CreatePlaceRequest true "Place"
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Success 201 "Place created"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /places [post]
func (h Handler) CreatePlace(c echo.Context) error {
	request := models.CreatePlaceRequest{}
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	err = h.usecase.CreatePlace(c.Request().Context(), request)
	if err != nil {
		return err
	}
	return c.NoContent(201)
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
	place, err := h.usecase.GetPlace(c.Request().Context(), id)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	return c.JSON(200, place)
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
// @Success 200 {array} models.Place
// @Failure 500 "Internal server error"
// @Router /places [get]
func (h Handler) GetPlaces(c echo.Context) error {
	request := models.GetPlacesRequest{}
	latitudeStr := c.QueryParam("latitude")
	longitudeStr := c.QueryParam("longitude")
	latitudeDeltaStr := c.QueryParam("latitudeDelta")
	longitudeDeltaStr := c.QueryParam("longitudeDelta")
	labels := c.QueryParam("labels")
	request.Labels = []string{labels}

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

	places, err := h.usecase.GetPlaces(c.Request().Context(), request)
	if err != nil {
		return err
	}
	return c.JSON(200, places)
}

// UpdatePlace godoc
// @Summary Update place
// @Description Update place
// @Tags places
// @Accept json
// @Produce json
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Param place body models.UpdatePlaceRequest true "Place"
// @Success 200 "Place updated"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /places [put]
func (h Handler) UpdatePlace(c echo.Context) error {
	request := models.UpdatePlaceRequest{}
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	err = h.usecase.UpdatePlace(c.Request().Context(), request)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}

// DeletePlace godoc
// @Summary Delete place
// @Description Delete place
// @Tags places
// @Accept json
// @Produce json
//
//	@Param			Authorization	header		string	true	"Authentication header"
//
// @Param id path string true "Place ID"
// @Success 200 "Place deleted"
// @Failure 404 "Place not found"
// @Failure 500 "Internal server error"
// @Router /places/{id} [delete]
func (h Handler) DeletePlace(c echo.Context) error {
	id := c.Param("id")
	err := h.usecase.DeletePlace(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}
