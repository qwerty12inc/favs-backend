package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e.GET("/health/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "Api is up and running!")
	})

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":" + serviceCfg.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("Server exiting...")

	return nil
}
