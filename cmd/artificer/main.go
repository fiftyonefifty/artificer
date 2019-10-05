package main

import (
	"artificer/pkg/api/handlers"
	"fmt"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config/appsettings.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func Alive() {
	now := time.Now().UTC()
	fmt.Println(fmt.Sprintf("Alive:%s", now))
}
func main() {
	// Creating a new Echo instance.
	e := echo.New()

	firstAlive := make(chan bool, 1)
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		Alive()
		firstAlive <- true
	})
	go func() {
		time.Sleep(5 * time.Second)
		Alive()
		firstAlive <- true
	}()

	c.Start()
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
	e.GET("/.well-known/openid-configuration", handlers.WellKnownOpenidConfiguration)
	e.GET("/.well-known/openid-configuration/jwks", handlers.WellKnownOpenidConfigurationJwks)
	// V1 Routes
	// v1 := e.Group("/v1")
	// v1Tokens := v1.Group("/tokens")
	// v1Tokens.GET("/tokens", handlers)
	fmt.Println()
	fmt.Println("------ Waiting for initial go routines to complete ------")
	fmt.Println()

	<-firstAlive
	// If you start me up, I'll never stop
	e.Logger.Fatal(e.Start(":8000"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
