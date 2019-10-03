package handlers

import (
	"artificer/pkg/api/renderings"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func WellKnownOpenidConfigurationJwks(c echo.Context) error {
	jwk := renderings.JwkResponse{}

	resp := renderings.WellKnownOpenidConfigurationJwksResponse{}
	resp.Keys = append(resp.Keys, jwk)

	return c.JSON(http.StatusOK, resp)
}
