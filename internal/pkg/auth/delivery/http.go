package delivery

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
	"strconv"
)

type Handler struct {
	// The Handler struct is the main entry point for the HTTP server. It contains a
	// reference to the UserUsecase, which is used to handle all business logic.
	usecase auth.Usecase
}

func NewHandler(usecase auth.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) SignUp(c echo.Context) error {
	request := models.SignUpRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}

	token, status := h.usecase.SignUp(c.Request().Context(), request)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, token)
}

func (h *Handler) Login(c echo.Context) error {
	request := models.LoginRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}

	token, status := h.usecase.Login(c.Request().Context(), request)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, token)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	request := models.UpdateUserRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}

	user, status := h.usecase.UpdateUser(c.Request().Context(), request)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, user)
}

func (h *Handler) CheckUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}

	user, status := h.usecase.GetUserByID(c.Request().Context(), id)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, user)
}

func (h *Handler) Logout(c echo.Context) error {
	token := c.Param("token")
	_, status := h.usecase.Logout(c.Request().Context(), token)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, "OK")
}

func (h *Handler) ActivateUser(c echo.Context) error {
	request := models.ActivateUserRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}

	status := h.usecase.ActivateUser(c.Request().Context(), request)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, "OK")
}

func (h *Handler) GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}
	user, status := h.usecase.GetUserByID(c.Request().Context(), id)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, user)
}
