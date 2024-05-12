package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type AuthRepositoryImpl struct {
	cl *firestore.Client
}

func NewAuthRepositoryImpl(cl *firestore.Client) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{cl: cl}
}

func (r *AuthRepositoryImpl) StoreToken(ctx context.Context, telegramID, token string) models.Status {
	_, err := r.cl.Collection("tokens").Doc(telegramID).Set(ctx,
		models.Token{
			Token: token,
		}, firestore.MergeAll)
	if err != nil {
		return models.Status{
			Code: models.InternalError,
		}
	}
	return models.Status{
		Code: models.OK,
	}
}

func (r *AuthRepositoryImpl) GetToken(ctx context.Context, telegramID string) (string, models.Status) {
	doc, err := r.cl.Collection("tokens").Doc(telegramID).Get(ctx)
	if err != nil {
		return "", models.Status{
			Code: models.NotFound,
		}
	}
	var token models.Token
	if err := doc.DataTo(&token); err != nil {
		return "", models.Status{
			Code: models.InternalError,
		}
	}
	return token.Token, models.Status{
		Code: models.OK,
	}
}
