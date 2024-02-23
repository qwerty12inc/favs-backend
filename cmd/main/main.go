package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth/delivery"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth/repository"
	"gitlab.com/v.rianov/favs-backend/internal/pkg/auth/usecase"
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
	if err != nil {
		return err
	}
	defer client.Close()

	smtpProvider, err := setupSMTP()
	if err != nil {
		return err
	}

	tokenProvider, err := setupTokenProvider()
	if err != nil {
		return err
	}

	activationCodesRepository, err := setupActivationCodesRepository()
	if err != nil {
		return err
	}

	repo := repository.NewFirestoreRepository(client)
	usecase := usecase.NewUsecase(repo, smtpProvider, tokenProvider, activationCodesRepository)
	handler := delivery.NewHandler(usecase)

	apiV1Group := e.Group("/api/v1")

	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.POST("/signup", handler.SignUp)
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/logout", handler.Logout)
	}

	userGroup := apiV1Group.Group("/user")
	{
		userGroup.GET("/me", handler.GetMe)
		userGroup.POST("/activation", handler.ActivateUser)
		userGroup.GET("/user/:id", handler.GetUserByID)
		userGroup.PUT("/user", handler.UpdateUser)
	}

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
