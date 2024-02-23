package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"strconv"
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
	_, _, err := r.cl.Collection("users").Add(ctx, user)
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
	doc.DataTo(&user)
	return user, models.Status{Code: models.OK}
}

func (r *FirestoreRepository) GetUserByID(ctx context.Context, id int) (models.User, models.Status) {
	doc, err := r.cl.Collection("users").Doc(strconv.Itoa(id)).Get(ctx)
	if err != nil {
		return models.User{}, models.Status{Code: models.NotFound, Message: err.Error()}
	}
	var user models.User
	doc.DataTo(&user)
	return user, models.Status{Code: models.OK}
}

func (r *FirestoreRepository) UpdateUser(ctx context.Context, user models.User) (models.User, models.Status) {
	_, err := r.cl.Collection("users").Doc(strconv.Itoa(user.ID)).Set(ctx, user)
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
