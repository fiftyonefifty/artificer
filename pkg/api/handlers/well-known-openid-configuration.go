package handlers

import (
	"artificer/pkg/api/renderings"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func WellKnownOpenidConfiguration(c echo.Context) error {
	resp := renderings.WellKnownOpenidConfigurationResponse{}
	return c.JSON(http.StatusOK, resp)
}
