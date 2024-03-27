package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "gitlab.com/v.rianov/favs-backend/docs"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	pkgmaps "gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	middleware2 "gitlab.com/v.rianov/favs-backend/internal/pkg/middleware"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/delivery"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/repository"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/usecase"
)

type ServiceConfig struct {
	Port string
}

const defaultPort = "8080"

func NewServiceConfig() ServiceConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return ServiceConfig{
		Port: port,
	}
}

// @title           Favs API
// @version         0.2.0
// @description     This is a documentation for favs API endpoints.

// @contact.name   API Maintainer
// @contact.email v.rianov@kabanov.agency

// @host 34.159.168.142
// @BasePath  /api/v1

// @securityDefinitions.basic  ApiKeyAuth
func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	serviceCfg := NewServiceConfig()
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := setupFirestore(ctx)
	log.Info("Firestore client created", err)
	if err != nil {
		return err
	}
	defer client.Close()

	authClient, err := setupFirebaseAuth(ctx)
	log.Info("Firebase auth client created", err)
	if err != nil {
		return err
	}

	authMiddleware := middleware2.NewAuthMiddlewareHandler(authClient)
	_ = authMiddleware

	apiV1Group := e.Group("/api/v1")

	apiV1Group.GET("/swagger/*", echoSwagger.WrapHandler)

	sheetsParser, err := setupSheetsParser(ctx)
	if err != nil {
		return err
	}

	cl, err := setupMapsClient()
	if err != nil {
		log.Info("Failed to create maps client ", err)
	}

	storageCLient, err := setupStorageClient(ctx)
	if err != nil {
		log.Info("Failed to create storage client", err)
		return err
	}

	storageRepo := repository.NewStorageRepository(storageCLient, os.Getenv("PLACES_BUCKET_ID"))

	// Place handlers
	placeRepo := repository.NewRepository(client)
	placeUsecase := usecase.NewUsecase(placeRepo,
		pkgmaps.NewLocationLinkResolver(cl), sheetsParser, storageRepo)
	placeHandler := delivery.NewHandler(placeUsecase)
	placeGroup := apiV1Group.Group("/places", authMiddleware.Auth)
	{
		placeGroup.GET("/:id", placeHandler.GetPlace)
		placeGroup.GET("", placeHandler.GetPlaces)
		placeGroup.GET("/:id/photos", placeHandler.GetPlacePhotos)
	}

	cityGroup := apiV1Group.Group("/cities", authMiddleware.Auth)
	{
		cityGroup.GET("", placeHandler.GetCities)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Info("Recovered from panic", r)
			}
		}()
		cities := []string{"Amsterdam", "Milan"}
		for _, city := range cities {
			status := placeUsecase.ImportPlacesFromSheet(ctx, fmt.Sprintf("%s!A2:G", city),
				city, "food", false)
			if status.Code != models.OK {
				log.Info("Failed to import places from sheet", status)
			} else {
				log.Info("Places imported from sheet", city)
			}
		}
	}()

	// Health check
	e.GET("/health/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "Api is up and running!")
	})

	sigCtx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":" + serviceCfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-sigCtx.Done()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("Server exiting...")

	return nil
}
