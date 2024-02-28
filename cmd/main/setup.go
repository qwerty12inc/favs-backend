package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"os"
)

func checkEnvVar(names ...string) error {
	for _, name := range names {
		if os.Getenv(name) == "" {
			return errors.New(name + " is not set")
		}
	}
	return nil
}

const serviceAccountPathEnv = "SERVICE_ACCOUNT_PATH"

func setupFirestore(ctx context.Context) (*firestore.Client, error) {
	err := checkEnvVar(serviceAccountPathEnv)
	if err != nil {
		return nil, err
	}
	sa := option.WithCredentialsFile(os.Getenv(serviceAccountPathEnv))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func setupFirebaseAuth(ctx context.Context) (*auth.Client, error) {
	err := checkEnvVar(serviceAccountPathEnv)
	if err != nil {
		return nil, err
	}
	sa := option.WithCredentialsFile(os.Getenv(serviceAccountPathEnv))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
