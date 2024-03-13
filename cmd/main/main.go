package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "gitlab.com/v.rianov/favs-backend/docs"
	"gitlab.com/v.rianov/favs-backend/internal/models"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	middleware2 "gitlab.com/v.rianov/favs-backend/internal/pkg/middleware"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/delivery"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/repository"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
// @version         0.1.0
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := setupFirestore(ctx)
	log.Println("Firestore client created", err)
	if err != nil {
		return err
	}
	defer client.Close()

	authClient, err := setupFirebaseAuth(ctx)
	log.Println("Firebase auth client created", err)
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

	// Place handlers
	placeRepo := repository.NewRepository(client)
	placeUsecase := usecase.NewUsecase(placeRepo, maps.LocationLinkResolverImpl{}, sheetsParser)
	placeHandler := delivery.NewHandler(placeUsecase)
	placeGroup := apiV1Group.Group("/places", authMiddleware.Auth)
	{
		placeGroup.POST("", placeHandler.CreatePlace)
		placeGroup.GET("/:id", placeHandler.GetPlace)
		placeGroup.GET("", placeHandler.GetPlaces)
		placeGroup.PUT("", placeHandler.UpdatePlace)
		placeGroup.DELETE("/:id", placeHandler.DeletePlace)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered from panic", r)
			}
		}()
		for {
			cities := []string{"Amsterdam", "Milan"}
			for _, city := range cities {
				status := placeUsecase.ImportPlacesFromSheet(ctx, fmt.Sprintf("%s!A2:G", city), city, false)
				if status.Code != models.OK {
					log.Println("Failed to import places from sheet", status)
				} else {
					log.Println("Places imported from sheet", city)
				}
			}
			time.Sleep(40 * time.Minute)
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
