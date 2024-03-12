package middleware

import (
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"log"
	"strings"
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
		log.Println("Token: ", token)
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			return c.JSON(401, "Unauthorized")
		}

		t, err := h.cl.VerifyIDToken(c.Request().Context(), token)
		if err != nil {
			return c.JSON(401, "Unauthorized")
		}

		firebaseUserInfo, err := h.cl.GetUser(c.Request().Context(), t.UID)
		if err != nil {
			return c.JSON(401, "Unauthorized")
		}
		log.Println("Firebase user email: ", firebaseUserInfo.UserInfo.Email)

		user := models.User{
			UID:   firebaseUserInfo.UID,
			Email: firebaseUserInfo.Email,
		}

		c.Set("user", user)
		c.Set("token", token)
		return next(c)
	}
}
