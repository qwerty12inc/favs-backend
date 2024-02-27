package usecase

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth"
	"log"
	"math/rand"
	"os"
)

type Usecase struct {
	repo          auth.Repository
	smtpProvider  auth.SMTPProvider
	tokenProvider auth.TokenProvider
	codeRepo      auth.ActivationCodesRepository
}

func NewUsecase(repo auth.Repository,
	smtpProvider auth.SMTPProvider,
	provider auth.TokenProvider,
	codeRepo auth.ActivationCodesRepository) *Usecase {
	return &Usecase{
		repo:          repo,
		smtpProvider:  smtpProvider,
		tokenProvider: provider,
		codeRepo:      codeRepo,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//
//func generateActivationLink(baseAddr, email, code string) string {
//	return baseAddr + "/api/v1/user/activation?email=" + email + "&code=" + code
//}

const baseAddrEnv = "BASE_ADDR"

func (u *Usecase) SignUp(ctx context.Context, request models.SignUpRequest) (string, models.Status) {
	_, status := u.repo.GetUserByEmail(ctx, request.Email)
	if status.Code == models.OK {
		return "", models.Status{Code: models.AlreadyExists, Message: "Email already registered"}
	}
	pswd := models.Password{}
	err := pswd.Set(request.Password)
	if err != nil {
		return "", models.Status{Code: models.InternalError, Message: "Failed to hash password"}
	}
	user, status := u.repo.SaveUser(ctx, models.User{
		ID:       uuid.New(),
		Email:    request.Email,
		Password: pswd,
	})
	if status.Code != models.OK {
		log.Println("failed to save user", status)
		return "", status
	}
	log.Println("user saved", user)

	token, status := u.tokenProvider.GenerateToken(ctx, user, false)
	if status.Code != models.OK {
		return "", status
	}
	log.Println("token generated", token)

	code := randSeq(12)
	status = u.codeRepo.SaveActivationCode(ctx, request.Email, code)
	if status.Code != models.OK {
		return "", status
	}
	log.Println("activation code saved", code)

	status = u.smtpProvider.Send(ctx, request.Email,
		"user_welcome.tmpl", map[string]string{"code": code})
	log.Println("email sent", status)
	if status.Code != models.OK {
		return "", status
	}

	return token, status
}

func (u *Usecase) Login(ctx context.Context, request models.LoginRequest) (string, models.Status) {
	user, status := u.repo.GetUserByEmail(ctx, request.Email)
	if status.Code != models.OK {
		return "", status
	}
	matches, err := user.Password.Matches(request.Password)
	if err != nil {
		return "", models.Status{Code: models.InternalError, Message: "Failed to compare passwords"}
	}
	if !matches {
		return "", models.Status{Code: models.InvalidCredentials, Message: "Invalid credentials"}
	}
	token, status := u.tokenProvider.GenerateToken(ctx, user, false)
	if status.Code != models.OK {
		return "", status
	}
	return token, status
}

func (u *Usecase) UpdateUser(ctx context.Context, request models.UpdateUserRequest) (models.User, models.Status) {
	user, status := u.repo.GetUserByID(ctx, request.ID)
	if status.Code != models.OK {
		return models.User{}, status
	}
	matches, err := user.Password.Matches(request.OldPassword)
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: "Failed to compare passwords"}
	}
	if !matches {
		return models.User{}, models.Status{Code: models.InvalidCredentials, Message: "Invalid credentials"}
	}
	pswd := models.Password{}
	err = pswd.Set(request.NewPassword)
	if err != nil {
		return models.User{}, models.Status{Code: models.InternalError, Message: "Failed to hash password"}
	}
	user.Password = pswd
	user, status = u.repo.UpdateUser(ctx, user)
	if status.Code != models.OK {
		return models.User{}, status
	}
	return user, status
}

func (u *Usecase) CheckUser(ctx context.Context, token string) (models.User, models.Status) {
	user, status := u.tokenProvider.ValidateToken(ctx, token)
	if status.Code != models.OK {
		return models.User{}, status
	}
	return user, status
}

func (u *Usecase) Logout(ctx context.Context, token string) (string, models.Status) {
	user, status := u.tokenProvider.ValidateToken(ctx, token)
	if status.Code != models.OK {
		return "", status
	}
	newToken, status := u.tokenProvider.GenerateToken(ctx, user, true)
	if status.Code != models.OK {
		return "", status
	}
	return newToken, status
}

func (u *Usecase) ActivateUser(ctx context.Context, request models.ActivateUserRequest) models.Status {
	storedCode, status := u.codeRepo.GetActivationCode(ctx, request.User.Email)
	log.Println("stored code", storedCode, status)
	if status.Code != models.OK {
		return status
	}
	if storedCode != request.Code {
		return models.Status{Code: models.InvalidToken, Message: "Invalid activation code"}
	}

	request.User.Activated = true
	_, status = u.repo.UpdateUser(ctx, request.User)
	if status.Code != models.OK {
		return status
	}
	return models.Status{Code: models.OK, Message: "User activated"}
}

func (u *Usecase) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, models.Status) {
	user, status := u.repo.GetUserByID(ctx, id)
	if status.Code != models.OK {
		return models.User{}, status
	}
	return user, status
}

func init() {
	baseAddr := os.Getenv(baseAddrEnv)
	if baseAddr == "" {
		panic("BASE_ADDR environment variable is not set")
	}
}
