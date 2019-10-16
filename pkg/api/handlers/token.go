package handlers

import (
	"artificer/pkg/api/renderings"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func TokenEndpoint(c echo.Context) (err error) {
	req := &TokenRequest{}
	if err = c.Bind(req); err != nil {
		return err
	}
	err = validateClient(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid request")
	}

	switch req.GrantType {
	case "client_credentials":
		return handleClientCredentialsFlow(c)
	case "arbitrary_resource_owner":
		return handleArbitraryResourceOwnerFlow(c)
	case "arbitrary_identity":

	}
	resp := renderings.HealthCheckResponse{
		Status: "Should Never See this",
	}
	return c.JSON(http.StatusOK, resp)
}
