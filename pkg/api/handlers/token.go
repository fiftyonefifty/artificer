package handlers

import (
	"artificer/pkg/api/renderings"
	"artificer/pkg/client/clientContext"
	"artificer/pkg/client/models"
	"artificer/pkg/userip"
	"context"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HealthCheck - Healthcheck Handler
func TokenEndpoint(c echo.Context) (err error) {
	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	req := &TokenRequest{}
	if err = c.Bind(req); err != nil {
		return err
	}
	// Store the user IP in ctx for use by code in other packages.
	userIP, err := userip.FromRequest(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctx = userip.NewContext(ctx, userIP)

	var client models.Client
	err, client = validateClient(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid request")
	}

	// a request has to be against the same client config.
	// i.e. a call to a database should happen once, because the client config can change in the database during the
	// lifetime of the request transaction.
	// We store the client in the contxt.Context and it can only be read from here going forward

	ctx = clientContext.NewContext(ctx, client)

	switch req.GrantType {
	case "client_credentials":
		return handleClientCredentialsFlow(ctx, c)
	case "arbitrary_resource_owner":
		return handleArbitraryResourceOwnerFlow(ctx, c)
	case "arbitrary_identity":

	}
	resp := renderings.HealthCheckResponse{
		Status: "Should Never See this",
	}
	return c.JSON(http.StatusOK, resp)
}
