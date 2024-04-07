package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go/v76"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/googlesheets"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"googlemaps.github.io/maps"
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

const placesBucketIDEnv = "PLACES_BUCKET_ID"

func setupStorageClient(ctx context.Context) (*storage.Client, error) {
	err := checkEnvVar(placesBucketIDEnv)
	if err != nil {
		return nil, err
	}
	err = checkEnvVar(serviceAccountPathEnv)
	if err != nil {
		return nil, err
	}
	sa := option.WithCredentialsFile(os.Getenv(serviceAccountPathEnv))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(ctx)
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

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := os.Getenv(tokenPath)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

const sheetsServiceAccountPathEnv = "SHEETS_SERVICE_ACCOUNT_PATH"
const spreadsheetIDEnv = "SPREADSHEET_ID"
const tokenPath = "SHEETS_TOKEN_PATH"

func setupSheetsParser(ctx context.Context) (*googlesheets.SheetsParser, error) {
	err := checkEnvVar(sheetsServiceAccountPathEnv, spreadsheetIDEnv, tokenPath)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(os.Getenv(sheetsServiceAccountPathEnv))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	if client == nil {
		log.Fatalf("Unable to retrieve client from config: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := os.Getenv(spreadsheetIDEnv)
	cl := googlesheets.NewSheetsParser(srv, spreadsheetId)
	return &cl, nil
}

const mapsAPIKeyEnv = "MAPS_API_KEY"

func setupMapsClient() (*maps.Client, error) {
	cl, err := maps.NewClient(maps.WithAPIKey(os.Getenv(mapsAPIKeyEnv)))
	if err != nil {
		log.Info("Failed to create maps client", err)
	}
	return cl, err
}

const stripeKeyEnv = "STRIPE_SECRET_KEY"

func setupStripe() {
	if os.Getenv(stripeKeyEnv) == "" {
		log.Fatal("Stripe key is not set")
	}
	stripe.Key = os.Getenv(stripeKeyEnv)
}
