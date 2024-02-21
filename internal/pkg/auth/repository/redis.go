package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"time"
)

type ActivatoinCodesRepository struct {
	Redis *redis.Client
}

func NewActivationCodesRepository(redis *redis.Client) *ActivatoinCodesRepository {
	return &ActivatoinCodesRepository{Redis: redis}
}

func (r *ActivatoinCodesRepository) SaveActivationCode(ctx context.Context, email, code string) models.Status {
	err := r.Redis.Set(ctx, email, code, time.Hour*3).Err()
	if err != nil {
		return models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return models.Status{Code: models.OK}
}

func (r *ActivatoinCodesRepository) GetActivationCode(ctx context.Context, email string) (string, models.Status) {
	code, err := r.Redis.Get(ctx, email).Result()
	if err != nil {
		return "", models.Status{Code: models.NotFound, Message: err.Error()}
	}
	return code, models.Status{Code: models.OK}
}
