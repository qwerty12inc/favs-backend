package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/redis/go-redis/v9"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth/repository"
	"os"
	"strconv"
)

const fireStoreProjectIDEnv = "FIRESTORE_PROJECT_ID"

func setupFirestore(ctx context.Context) (*firestore.Client, error) {
	client, err := firestore.NewClient(ctx, os.Getenv(fireStoreProjectIDEnv))
	return client, err
}

const (
	smtpHostEnv     = "SMTP_HOST"
	smtpPortEnv     = "SMTP_PORT"
	smtpUserEnv     = "SMTP_USER"
	smtpPasswordEnv = "SMTP_PASSWORD"
	smtpSenderEnv   = "SMTP_SENDER"
)

func setupSMTP() (*repository.Mailer, error) {
	host := os.Getenv(smtpHostEnv)
	portStr := os.Getenv(smtpPortEnv)
	user := os.Getenv(smtpUserEnv)
	password := os.Getenv(smtpPasswordEnv)
	sender := os.Getenv(smtpSenderEnv)

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return repository.NewMailer(host, port, user, password, sender), nil
}

const signingKeyEnv = "SIGNING_KEY"

func setupTokenProvider() (*repository.TokenProvider, error) {
	return repository.NewTokenProvider(os.Getenv(signingKeyEnv)), nil
}

const (
	redisAddrEnv = "REDIS_ADDR"
	redisPassEnv = "REDIS_PASS"
	redisUserEnv = "REDIS_USER"
)

func setupActivationCodesRepository() (*repository.ActivatoinCodesRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(redisAddrEnv),
		Password: os.Getenv(redisPassEnv),
		Username: os.Getenv(redisUserEnv),
	})
	return repository.NewActivationCodesRepository(client), nil
}
