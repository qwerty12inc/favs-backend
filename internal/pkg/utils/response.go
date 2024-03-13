package utils

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

func HandleResponse(c echo.Context, status models.Status, body interface{}) error {
	switch status.Code {
	case models.OK:
		return c.JSON(200, body)
	case models.BadRequest:
		return c.JSON(400, status)
	case models.NotFound:
		return c.JSON(404, status)
	case models.InternalError:
		return c.JSON(500, status)
	case models.Unauthorized:
		return c.JSON(401, status)
	case models.Forbidden:
		return c.JSON(403, status)
	default:
		return c.JSON(500, models.Status{
			Code:    models.InternalError,
			Message: "Internal server error",
		})
	}
}
