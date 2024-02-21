package models

type StatusCode int

const (
	OK StatusCode = iota
	InvalidRequest
	InvalidCredentials
	InvalidToken
	InternalError
	AlreadyExists
	NotFound
	BadRequest
	Unauthorized
	Forbidden
)

type Status struct {
	Code    StatusCode `json:"code"`
	Message string     `json:"message"`
}
