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
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
	Email       string `json:"email"`
}

type ActivateUserRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
