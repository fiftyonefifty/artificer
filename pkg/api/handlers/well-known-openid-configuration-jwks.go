package handlers

import (
	"artificer/pkg/api/renderings"
	"artificer/pkg/keyvault"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	gocache "github.com/pmylund/go-cache"
)

var (
	cache    = gocache.New(24*time.Hour, time.Hour)
	cacheKey = "85b75fb0-f120-4bfb-a0fe-f017cc72e41f"
)

func DoKeyvaultBackground() (err error) {
	now := time.Now().UTC()
	fmt.Println(fmt.Sprintf("Start-DoKeyvaultBackground:%s", now))
	ctx := context.Background()
	activeKeys, _, err := keyvault.GetActiveKeysVersion(ctx)
	if err != nil {
		return
	}
	resp := renderings.WellKnownOpenidConfigurationJwksResponse{}

	for _, element := range activeKeys {

		jwk := renderings.JwkResponse{}
		jwk.Kid = *element.Key.Kid
		jwk.Kty = string(element.Key.Kty)
		jwk.N = *element.Key.N
		jwk.E = *element.Key.E
		jwk.Alg = "RSA256"
		jwk.Use = "sig"
		resp.Keys = append(resp.Keys, jwk)
	}
	cache.Set(cacheKey, resp, gocache.NoExpiration)
	fmt.Println(fmt.Sprintf("Success-DoKeyvaultBackground:%s", now))
	return
}

// HealthCheck - Healthcheck Handler
func WellKnownOpenidConfigurationJwks(c echo.Context) error {

	cachedResponse, found := cache.Get(cacheKey)
	if !found {
		err := DoKeyvaultBackground()
		if err != nil {
			log.Fatalf("failed to DoKeyvaultBackground: %v\n", err.Error())
			return c.JSON(http.StatusBadRequest, nil)
		}
		cachedResponse, found = cache.Get(cacheKey)
		if !found {
			log.Fatalf("critical failure to DoKeyvaultBackground:\n")
			return c.JSON(http.StatusBadRequest, nil)
		}
	}

	return c.JSON(http.StatusOK, cachedResponse)
}
