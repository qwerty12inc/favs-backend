package models

import "github.com/google/uuid"

type SignUpRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	ID          uuid.UUID `json:"id,omitempty"`
	NewPassword string    `json:"new_password,omitempty"`
	OldPassword string    `json:"old_password,omitempty"`
	Email       string    `json:"email,omitempty"`
}

type ActivateUserRequest struct {
	Code string `json:"code"`
	User User   `json:"user"`
}
