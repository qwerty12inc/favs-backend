package middleware

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
)

type AuthMiddlewareHandler struct {
	tokenProvider auth.TokenProvider
	repository    auth.Repository
}

func (h AuthMiddlewareHandler) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(401, "Unauthorized")
		}
		user, status := h.tokenProvider.ValidateToken(c.Request().Context(), token)
		if status.Code != models.OK {
			return c.JSON(401, "Unauthorized")
		}
		user, status = h.repository.GetUserByID(c.Request().Context(), user.ID)
		if status.Code != models.OK {
			return c.JSON(401, "Unauthorized")
		}
		c.Set("user", user)
		return next(c)
	}
}
