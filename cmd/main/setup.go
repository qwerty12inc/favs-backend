package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth/repository"
	"os"
	"strconv"
)

func checkEnvVar(names ...string) error {
	for _, name := range names {
		if os.Getenv(name) == "" {
			return errors.New(name + " is not set")
		}
	}
	return nil
}

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
	err := checkEnvVar(smtpHostEnv, smtpPortEnv, smtpUserEnv, smtpPasswordEnv, smtpSenderEnv)
	if err != nil {
		return nil, err
	}

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
	err := checkEnvVar(signingKeyEnv)
	if err != nil {
		return nil, err
	}
	return repository.NewTokenProvider(os.Getenv(signingKeyEnv)), nil
}

const (
	redisAddrEnv = "REDIS_ADDR"
	redisPassEnv = "REDIS_PASS"
)

func setupActivationCodesRepository() (*repository.ActivatoinCodesRepository, error) {
	err := checkEnvVar(redisAddrEnv, redisPassEnv)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(redisAddrEnv),
		Password: os.Getenv(redisPassEnv),
	})
	return repository.NewActivationCodesRepository(client), nil
}
