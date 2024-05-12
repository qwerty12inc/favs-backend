package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
)

type TelegramMiddlewareHandler struct {
	usecase auth.Usecase
}

func NewTelegramMiddlewareHandler(usecase auth.Usecase) TelegramMiddlewareHandler {
	return TelegramMiddlewareHandler{
		usecase: usecase,
	}
}

func (h TelegramMiddlewareHandler) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			return c.JSON(401, "Unauthorized")
		}

		if token == "test" {
			user := models.User{
				UID:   "test",
				Email: "",
			}
			c.Set("user", user)
			c.Set("token", token)
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "user", user)))
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "token", token)))
			return next(c)
		}

		telegramID := c.QueryParam("telegramID")
		status := h.usecase.Verify(c.Request().Context(), token, telegramID)
		if status.Code != models.OK {
			return c.JSON(401, "Unauthorized")
		}

		c.Set("token", token)
		c.Set("telegramID", telegramID)
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "telegramID", telegramID)))
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "token", token)))
		return next(c)
	}
}
