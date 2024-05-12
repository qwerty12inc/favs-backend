package usecase

import (
	"context"
	"math/rand"

	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
)

type UsecaseImpl struct {
	repo auth.Repository
}

func NewUsecase(repo auth.Repository) UsecaseImpl {
	return UsecaseImpl{repo: repo}
}

func generateToken() string {
	token := ""
	for i := 0; i < 32; i++ {
		token += string(rune(rand.Intn(26) + 65))
	}
	return token
}

func (u UsecaseImpl) Login(ctx context.Context, telegramID string) (string, models.Status) {
	token := generateToken()
	status := u.repo.StoreToken(ctx, telegramID, token)
	return token, status
}

func (u UsecaseImpl) Verify(ctx context.Context, token, telegramID string) models.Status {
	existingToken, status := u.repo.GetToken(ctx, telegramID)
	if status.Code != models.OK {
		return models.Status{
			Code: models.NotFound,
		}
	}
	if existingToken != token {
		return models.Status{
			Code: models.Unauthorized,
		}
	}
	return models.Status{
		Code: models.OK,
	}
}
