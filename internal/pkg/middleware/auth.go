package middleware

import (
	"context"
	"os"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type AuthMiddlewareHandler struct {
	cl *auth.Client
}

func NewAuthMiddlewareHandler(cl *auth.Client) AuthMiddlewareHandler {
	return AuthMiddlewareHandler{
		cl: cl,
	}
}

func (h AuthMiddlewareHandler) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		log.Debug("Token: ", token)
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			return c.JSON(401, "Unauthorized")
		}

		log.Debug("Service token: ", os.Getenv("SERVICE_TOKEN"))

		if token == os.Getenv("SERVICE_TOKEN") {
			user := models.User{
				UID:   "openapp",
				Email: "openapp@openapp.com",
			}
			c.Set("user", user)
			c.Set("token", token)
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "user", user)))
			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "token", token)))
			return next(c)
		}

		t, err := h.cl.VerifyIDToken(c.Request().Context(), token)
		if err != nil {
			return c.JSON(401, "Unauthorized")
		}

		firebaseUserInfo, err := h.cl.GetUser(c.Request().Context(), t.UID)
		if err != nil {
			return c.JSON(401, "Unauthorized")
		}
		log.Debug("Firebase user email: ", firebaseUserInfo.UserInfo.Email)

		user := models.User{
			UID:   firebaseUserInfo.UID,
			Email: firebaseUserInfo.Email,
		}

		c.Set("user", user)
		c.Set("token", token)
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "user", user)))
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "token", token)))
		return next(c)
	}
}
