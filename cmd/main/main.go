package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
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

	return e.Start(":" + serviceCfg.Port)
}
