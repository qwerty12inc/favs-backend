package delivery

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/utils"
)

type AuthHandler struct {
	uc auth.Usecase
}

func NewAuthHandler(uc auth.Usecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param telegramID query string true "Telegram ID"
// @Success 200 {string} token
// @Failure 400 "Invalid request"
// @Failure 500 "Internal server error"
// @Router /tg/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	telegramID := c.QueryParam("telegramID")
	token, status := h.uc.Login(c.Request().Context(), telegramID)
	return utils.HandleResponse(c, status, token)
}

// Verify godoc
// @Summary Verify
// @Description Verify
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Token"
// @Param telegramID query string true "Telegram ID"
// @Success 200 "OK"
// @Failure 400 "Invalid request"
// @Failure 500 "Internal server error"
// @Router /tg/verify [post]
func (h *AuthHandler) Verify(c echo.Context) error {
	token := c.QueryParam("token")
	telegramID := c.QueryParam("telegramID")
	status := h.uc.Verify(c.Request().Context(), token, telegramID)
	return utils.HandleResponse(c, status, nil)
}
