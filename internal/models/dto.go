package models

type SignUpRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UpdateUserRequest struct {
	ID          int    `json:"id,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	Email       string `json:"email,omitempty"`
}

type ActivateUserRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
