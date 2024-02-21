package delivery

import (
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
)

type Handler struct {
	// The Handler struct is the main entry point for the HTTP server. It contains a
	// reference to the UserUsecase, which is used to handle all business logic.
	usecase auth.Usecase
}

func NewHandler(usecase auth.Usecase) *Handler {
	return &Handler{usecase: usecase}
}
