package handlers

import (
	"artificer/pkg/appError"
	"artificer/pkg/client/clientContext"
	"artificer/pkg/client/models"
	"artificer/pkg/userip"
	"context"
	"log"
	"net/http"
	"time"

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
	ctx = context.Background()

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

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	var client models.Client
	err, client = validateClient(ctx, req)
	if err != nil {
		log.Printf("Failed to validate client:%s\n", err.Error())
		switch err.(type) {
		case *appError.AppError:
			p := err.(*appError.AppError)
			return c.JSON(p.Code, p.Message)
		default:
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
	}

	// a request has to be against the same client config.
	// i.e. a call to a database should happen once, because the client config can change in the database during the
	// lifetime of the request transaction.
	// We store the client in the contxt.Context and it can only be read from here going forward

	ctx = clientContext.NewContext(ctx, client)

	switch req.GrantType {
	case "client_credentials":
		err = handleClientCredentialsFlow(ctx, c)
	case "arbitrary_resource_owner":
		err = handleArbitraryResourceOwnerFlow(ctx, c)
	default:
		err = appError.New(http.StatusBadRequest, "should-never-see-this")
	}
	if err != nil {
		log.Printf("Failed grant_type:%s\n", req.GrantType)
		switch err.(type) {
		case *appError.AppError:
			p := err.(*appError.AppError)
			return c.JSON(p.Code, p.Message)
		default:
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
	} else {
		return
	}

}
