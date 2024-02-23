package repository

import (
	"context"
	"database/sql"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

const (
	InsertUser        = `INSERT INTO users (email, password, is_active) VALUES ($1, $2, $3) RETURNING id`
	SelectUserByEmail = `SELECT id, email, password, is_active, created_at FROM users WHERE email = $1`
	SelectUserByID    = `SELECT id, email, password, is_active, created_at FROM users WHERE id = $1`
	UpdateUser        = `UPDATE users SET email = $1, password = $2, is_active = $3 WHERE id = $4`
	DeleteUser        = `DELETE FROM users WHERE id = $1`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveUser(ctx context.Context, user models.User) (models.User, models.Status) {
	res, err := r.db.ExecContext(ctx, InsertUser, user.Email, user.Password, user.Activated)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.User{}, models.Status{Code: models.AlreadyExists, Message: "user with this email already exists"}
		}
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	user.ID = int(id)
	return user, models.Status{Code: models.OK}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (models.User, models.Status) {
	var user models.User
	err := r.db.QueryRowContext(ctx, SelectUserByEmail, email).Scan(&user.ID, &user.Email, &user.Password,
		&user.Activated, &user.CreatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.User{}, models.Status{Code: models.NotFound, Message: "user with this email not found"}
		}
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return user, models.Status{Code: models.OK}
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (models.User, models.Status) {
	var user models.User
	err := r.db.QueryRowContext(ctx, SelectUserByID, id).Scan(&user.ID, &user.Email, &user.Password,
		&user.Activated, &user.CreatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.User{}, models.Status{Code: models.NotFound, Message: "user with this id not found"}
		}
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return user, models.Status{Code: models.OK}
}

func (r *Repository) UpdateUser(ctx context.Context, user models.User) (models.User, models.Status) {
	_, err := r.db.ExecContext(ctx, UpdateUser, user.Email, user.Password, user.Activated, user.ID)
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return user, models.Status{Code: models.OK}
}

func (r *Repository) DeleteUser(ctx context.Context, id string) models.Status {
	_, err := r.db.ExecContext(ctx, DeleteUser, id)
	if err != nil {
		return models.Status{Code: models.InternalError, Message: err.Error()}
	}
	return models.Status{Code: models.OK}
}
