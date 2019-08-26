package handlers

import (
	"artificer/pkg/api/renderings"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func HealthCheck(c echo.Context) error {
	resp := renderings.HealthCheckResponse{
		Status: "UP",
	}
	return c.JSON(http.StatusOK, resp)
}
