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
	//E := viper.GetString("keyVault.clientId")

	ctx := context.Background()
	//defer resources.Cleanup(ctx)

	activeKeys, err := keyvault.GetActiveKeysVersion(ctx)
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

	return c.JSON(http.StatusOK, resp)
}
