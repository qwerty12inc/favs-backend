package usecase

import (
	"context"
	"math/rand"

	"github.com/labstack/gommon/log"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
)

type AuthUsecaseImpl struct {
	repo auth.Repository
}

func NewAuthUsecaseImpl(repo auth.Repository) AuthUsecaseImpl {
	return AuthUsecaseImpl{repo: repo}
}

func generateToken() string {
	token := ""
	for i := 0; i < 32; i++ {
		token += string(rune(rand.Intn(26) + 65))
	}
	return token
}

func (u AuthUsecaseImpl) Login(ctx context.Context, telegramID string) (string, models.Status) {
	token := generateToken()
	log.Info("Generated token: ", token)
	status := u.repo.StoreToken(ctx, telegramID, token)
	if status.Code != models.OK {
		log.Error("Failed to store token ", status)
		return "", status
	}
	return token, status
}

func (u AuthUsecaseImpl) Verify(ctx context.Context, token, telegramID string) models.Status {
	existingToken, status := u.repo.GetToken(ctx, telegramID)
	if status.Code != models.OK {
		log.Error("Failed to get token ", status)
		return models.Status{
			Code: models.NotFound,
		}
	}
	if existingToken != token {
		log.Error("Token mismatch")
		return models.Status{
			Code: models.Unauthorized,
		}
	}
	log.Info("Token verified")
	return models.Status{
		Code: models.OK,
	}
}
