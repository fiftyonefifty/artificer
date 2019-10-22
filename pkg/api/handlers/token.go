package handlers

import (
	"artificer/pkg/api/renderings"
	"artificer/pkg/client/models"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func TokenEndpoint(c echo.Context) (err error) {
	req := &TokenRequest{}
	if err = c.Bind(req); err != nil {
		return err
	}
	var client models.Client
	err, client = validateClient(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid request")
	}

	// a request has to be against the same client config.
	// i.e. a call to a database should happen once, because the client config can change in the database during the
	// lifetime of the request transaction.
	// We store the client in the echo.Context and it can only be read from here going forward
	c.Set("_client", client)

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
