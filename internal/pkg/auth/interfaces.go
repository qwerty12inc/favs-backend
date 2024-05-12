package auth

import (
	"context"

	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type Usecase interface {
	Login(ctx context.Context, telegramID string) (string, models.Status)
	Verify(ctx context.Context, token, telegramID string) models.Status
}

type Repository interface {
	StoreToken(ctx context.Context, telegramID, token string) models.Status
	GetToken(ctx context.Context, telegramID string) (string, models.Status)
}
