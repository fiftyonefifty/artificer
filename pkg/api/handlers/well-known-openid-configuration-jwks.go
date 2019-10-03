package handlers

import (
	"artificer/pkg/api/renderings"
	"context"
	"log"
	"net/http"

	"artificer/pkg/config"
	"artificer/pkg/keyvault"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func WellKnownOpenidConfigurationJwks(c echo.Context) error {

	var err error
	err = config.ParseEnvironment()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err.Error())
	}

	ctx := context.Background()
	//defer resources.Cleanup(ctx)
	a, err := keyvault.GetKeysVersion(ctx)

	keyItem := a.Values()[0]

	jwk := renderings.JwkResponse{}
	jwk.Kty = config.ClientID()
	jwk.Kid = *keyItem.Kid
	resp := renderings.WellKnownOpenidConfigurationJwksResponse{}
	resp.Keys = append(resp.Keys, jwk)

	return c.JSON(http.StatusOK, resp)
}
