package delivery

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/usecase"
)

type Handler struct {
	usecase usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

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

func (h Handler) GetPlace(c echo.Context) error {
	id := c.Param("id")
	place, err := h.usecase.GetPlace(c.Request().Context(), id)
	if err != nil {
		c.String(500, err.Error())
		return err
	}
	return c.JSON(200, place)
}

func (h Handler) GetPlaces(c echo.Context) error {
	request := models.GetPlacesRequest{}
	err := c.Bind(&request)
	if err != nil {
		return err
	}
	places, err := h.usecase.GetPlaces(c.Request().Context(), request)
	if err != nil {
		return err
	}
	return c.JSON(200, places)
}

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

func (h Handler) DeletePlace(c echo.Context) error {
	id := c.Param("id")
	err := h.usecase.DeletePlace(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}
