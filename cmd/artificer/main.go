package main

import (
	"artificer/pkg/api/handlers"
	"artificer/pkg/config"
	"artificer/pkg/keyvault"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

var (
	ProcessDirectory string
)

func init() {
	var err error
	ProcessDirectory, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ProcessDirectory)

	viper.SetConfigFile("config/appsettings.json")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	config.LoadClientConfig(ProcessDirectory)
}

func Alive() {
	now := time.Now().UTC()
	fmt.Println(fmt.Sprintf("Alive:%s", now))
}
func main() {
	var err error

	err = config.ParseEnvironment()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err.Error())
	}

	// Creating a new Echo instance.
	e := echo.New()

	firstAlive := make(chan bool, 1)

	go func() {
		keyvault.DoKeyvaultBackground()
		firstAlive <- true
	}()
	go func() {
		ctx := context.Background()
		keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl")
		keyvault.CreateKey(ctx, keyVaultUrl, "test")

	}()

	c := cron.New()
	cronSpec := viper.GetString("keyVault.cronSpec") // i.e. "@every 10s"

	c.AddFunc(cronSpec, func() {
		keyvault.DoKeyvaultBackground()
		firstAlive <- true
	})
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
	e.GET("/mint-test-token", handlers.MintTestToken)
	e.POST("/connect/token", handlers.TokenEndpoint)

	// V1 Routes
	// v1 := e.Group("/v1")
	// v1Tokens := v1.Group("/tokens")
	// v1Tokens.GET("/tokens", handlers)
	fmt.Println()
	fmt.Println("------ Waiting for initial go routines to complete ------")
	fmt.Println()

	<-firstAlive
	// If you start me up, I'll never stop
	e.Logger.Fatal(e.Start(":9000"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
