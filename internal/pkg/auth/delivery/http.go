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

// SignUp godoc
// @Summary Sign up
// @Description Create a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.SignUpRequest true "User data"
// @Success 200
// @Failure 400 {string} string
// @Router /auth/signup [post]
func (h *Handler) SignUp(c echo.Context) error {
	request := models.SignUpRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}

	token, status := h.usecase.SignUp(c.Request().Context(), request)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}

	c.Response().Header().Set("Authorization", "Bearer "+token)
	return c.JSON(200, "OK")
}

// Login godoc
// @Summary Login
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "User data"
// @Success 200
// @Failure 400 {string} string
// @Router /auth/login [post]
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

// UpdateUser godoc
// @Summary Update user
// @Description Update user data
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UpdateUserRequest true "User data"
// @Success 200 {object} models.User
// @Failure 400 {string} string
// @Router /user [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	request := models.UpdateUserRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: err.Error()})
	}
	user, ok := c.Get("user").(models.User)
	if !ok {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: "User not found"})
	}
	request.ID = user.ID

	user, status := h.usecase.UpdateUser(c.Request().Context(), request)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, user)
}

// GetMe godoc
// @Summary Get user
// @Description Get user data
// @Tags user
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {string} string
// @Router /user/me [get]
func (h *Handler) GetMe(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: "User not found"})
	}
	return c.JSON(200, user)
}

// Logout godoc
// @Summary Logout
// @Description Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200
// @Failure 400 {string} string
// @Router /auth/logout/{token} [post]
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
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: "Code is empty"})
	}
	request.Code = code

	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(400, models.Status{Code: models.BadRequest, Message: "Email is empty"})
	}
	request.Email = email

	status := h.usecase.ActivateUser(c.Request().Context(), request)
	if status.Code != models.OK {
		return c.JSON(400, status)
	}
	return c.JSON(200, "OK")
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user data by ID
// @Tags user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string
// @Router /user/{id} [get]
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
