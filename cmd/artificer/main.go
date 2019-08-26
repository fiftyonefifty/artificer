package main

import (
	"artificer/pkg/api/handlers"
	"net/http"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Creating a new Echo instance.
	e := echo.New()

	// Configure Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	e.File("/", "static/index.html")

	// in order to serve static assets
	e.Static("/static", "static")

	// Route / to handler function
	e.GET("/health", handlers.HealthCheck)

	// V1 Routes
	// v1 := e.Group("/v1")
	// v1Tokens := v1.Group("/tokens")
	// v1Tokens.GET("/tokens", handlers)
	// Start Server
	e.Logger.Fatal(e.Start(":8000"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
