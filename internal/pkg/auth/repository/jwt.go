package repository

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/v.rianov/favs-backend/internal/models"
)

type TokenProvider struct {
	SigningKey string
}

type TokenClaims struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
}

func NewTokenProvider(signingKey string) *TokenProvider {
	return &TokenProvider{SigningKey: signingKey}
}

func (tp *TokenProvider) GenerateToken(ctx context.Context,
	user models.User, expiry bool) (string, models.Status) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID": user.ID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(tp.SigningKey))
	if err != nil {
		return "", models.Status{Code: models.InternalError, Message: err.Error()}
	}

	return tokenString, models.Status{Code: models.OK}
}

func (tp *TokenProvider) ValidateToken(ctx context.Context, tokenStr string) (models.User, models.Status) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(tp.SigningKey), nil
	})

	user := models.User{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.ID = claims["UserID"].(int)
	} else {
		return models.User{}, models.Status{Code: models.Unauthorized, Message: err.Error()}
	}

	return user, models.Status{Code: models.OK}
}
