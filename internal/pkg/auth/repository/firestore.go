package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type RepositoryImpl struct {
	cl *firestore.Client
}

func NewRepositoryImpl(cl *firestore.Client) *RepositoryImpl {
	return &RepositoryImpl{cl: cl}
}

func (r *RepositoryImpl) StoreToken(ctx context.Context, telegramID, token string) error {
	_, err := r.cl.Collection("tokens").Doc(telegramID).Set(ctx,
		models.Token{
			Token: token,
		}, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryImpl) GetToken(ctx context.Context, telegramID string) (string, error) {
	doc, err := r.cl.Collection("tokens").Doc(telegramID).Get(ctx)
	if err != nil {
		return "", err
	}
	var token models.Token
	if err := doc.DataTo(&token); err != nil {
		return "", err
	}
	return token.Token, nil
}
