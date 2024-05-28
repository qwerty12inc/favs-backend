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
	delivery2 "gitlab.com/v.rianov/favs-backend/internal/pkg/auth/delivery"
	repository2 "gitlab.com/v.rianov/favs-backend/internal/pkg/auth/repository"
	usecase2 "gitlab.com/v.rianov/favs-backend/internal/pkg/auth/usecase"
	pkgmaps "gitlab.com/v.rianov/favs-backend/internal/pkg/maps"
	middleware2 "gitlab.com/v.rianov/favs-backend/internal/pkg/middleware"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/delivery"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/repository"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/places/usecase"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/stripe"
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

	log.SetLevel(log.DEBUG)

	e.Use(middleware.CORS())
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

	setupStripe()
	stripeConnector := stripe.NewStripeConnector()

	// Place handlers
	placeRepo := repository.NewRepository(client)
	placeUsecase := usecase.NewUsecase(placeRepo,
		pkgmaps.NewLocationLinkResolver(cl), sheetsParser, storageRepo,
		stripeConnector)
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

	purchaseGroup := apiV1Group.Group("/purchases", authMiddleware.Auth)
	{
		purchaseGroup.POST("", placeHandler.SaveUserPurchase)
	}

	paymentLinkGroup := apiV1Group.Group("/payments", authMiddleware.Auth)
	{
		paymentLinkGroup.GET("", placeHandler.GeneratePaymentLink)
	}

	tgAuthRepo := repository2.NewAuthRepositoryImpl(client)
	tgUsecase := usecase2.NewAuthUsecaseImpl(tgAuthRepo)
	tgHandler := delivery2.NewAuthHandler(tgUsecase)

	tgGroup := apiV1Group.Group("/tg")
	{
		tgGroup.POST("/login", tgHandler.Login)
	}

	tgMiddleware := middleware2.NewTelegramMiddlewareHandler(tgUsecase)

	tgPlaceGroup := apiV1Group.Group("/tg/places", tgMiddleware.Auth)
	{
		// add OPTIONS method for CORS preflight requests
		tgPlaceGroup.GET("", placeHandler.TelegramGetPlaces)
		tgPlaceGroup.GET("/:id", placeHandler.TelegramGetPlace)
		tgPlaceGroup.GET("/:id/photos", placeHandler.GetPlacePhotos)
		tgPlaceGroup.POST("/:id/reports", placeHandler.SaveReport)
	}
	tgPlaceGroup.GET("/cities", placeHandler.TelegramGetCities, tgMiddleware.Auth)

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
