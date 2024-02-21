package auth

import (
	"context"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type Usecase interface {
	SignUp(ctx context.Context, request models.SignUpRequest) (string, models.Status)
	Login(ctx context.Context, request models.LoginRequest) (string, models.Status)
	UpdateUser(ctx context.Context, request models.UpdateUserRequest) (models.User, models.Status)
	CheckUser(ctx context.Context, token string) (models.User, models.Status)
	Logout(ctx context.Context, token string) (string, models.Status)
	ActivateUser(ctx context.Context, request models.ActivateUserRequest) models.Status
	GetUserByID(ctx context.Context, id int) (models.User, models.Status)
}

type Repository interface {
	SaveUser(ctx context.Context, user models.User) (models.User, models.Status)
	GetUserByEmail(ctx context.Context, email string) (models.User, models.Status)
	GetUserByID(ctx context.Context, id int) (models.User, models.Status)
	UpdateUser(ctx context.Context, user models.User) (models.User, models.Status)
}

type SMTPProvider interface {
	Send(ctx context.Context, recipient, templateFile string, data interface{}) models.Status
}

type TokenProvider interface {
	GenerateToken(ctx context.Context, user models.User, expiry bool) (string, models.Status)
	ValidateToken(ctx context.Context, token string) (models.User, models.Status)
}

type ActivationCodesRepository interface {
	SaveActivationCode(ctx context.Context, email, code string) models.Status
	GetActivationCode(ctx context.Context, email string) (string, models.Status)
}
