package middleware

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
	"log"
	"strings"
)

type AuthMiddlewareHandler struct {
	TokenProvider auth.TokenProvider
	Repository    auth.Repository
}

func (h AuthMiddlewareHandler) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		log.Println("Token: ", token)
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			return c.JSON(401, "Unauthorized")
		}
		user, status := h.TokenProvider.ValidateToken(c.Request().Context(), token)
		log.Println("User: ", user, "Status: ", status)
		if status.Code != models.OK {
			return c.JSON(401, "Unauthorized")
		}
		log.Println("User: ", user)
		user, status = h.Repository.GetUserByID(c.Request().Context(), user.ID)
		if status.Code != models.OK {
			return c.JSON(401, "Unauthorized")
		}
		c.Set("user", user)
		return next(c)
	}
}
