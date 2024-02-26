package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/google/uuid"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type FirestoreRepository struct {
	cl *firestore.Client
}

func NewFirestoreRepository(cl *firestore.Client) *FirestoreRepository {
	return &FirestoreRepository{
		cl: cl,
	}
}

func (r *FirestoreRepository) SaveUser(ctx context.Context, user models.User) (models.User, models.Status) {
	_, err := r.cl.Collection("users").Doc(user.ID.String()).Set(ctx, user)
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return user, models.Status{Code: models.OK}
}

func (r *FirestoreRepository) GetUserByEmail(ctx context.Context, email string) (models.User, models.Status) {
	iter := r.cl.Collection("users").Where("email", "==", email).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return models.User{}, models.Status{Code: models.NotFound, Message: err.Error()}
	}
	var user models.User
	err = doc.DataTo(&user)
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return user, models.Status{Code: models.OK}
}

func (r *FirestoreRepository) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, models.Status) {
	iter := r.cl.Collection("users").Where("id", "==", id).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return models.User{}, models.Status{Code: models.NotFound, Message: err.Error()}
	}
	var user models.User
	err = doc.DataTo(&user)
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return user, models.Status{Code: models.OK}
}

func (r *FirestoreRepository) UpdateUser(ctx context.Context, user models.User) (models.User, models.Status) {
	_, err := r.cl.Collection("users").Doc(user.ID.String()).Set(ctx, user)
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return user, models.Status{Code: models.OK}
}

func (r *FirestoreRepository) DeleteUser(ctx context.Context, id string) models.Status {
	_, err := r.cl.Collection("users").Doc(id).Delete(ctx)
	if err != nil {
		return models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return models.Status{Code: models.OK}
}
