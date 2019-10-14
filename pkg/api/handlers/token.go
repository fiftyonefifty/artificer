package handlers

import (
	"artificer/pkg/api/models"
	"artificer/pkg/api/renderings"
	"artificer/pkg/config"
	"artificer/pkg/util"
	"errors"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func validateClient(req *TokenRequest) (err error) {

	if len(req.ClientID) == 0 || len(req.ClientSecret) == 0 {
		err = errors.New("client_id or client_secret is not present")
		fmt.Println(err.Error())
		return
	}

	sEnc := util.StringSha256Encode64(req.ClientSecret)

	var client *models.Client
	client = config.ClientMap[req.ClientID]
	if client == nil {
		err = errors.New(fmt.Sprintf("client_id: %s does not exist", req.ClientID))
		fmt.Println(err.Error())
		return
	}

	foundSecret := false
	for _, element := range client.ClientSecrets {
		foundSecret = (sEnc == element.Value)
		if foundSecret {
			break
		}
	}
	if !foundSecret {
		err = errors.New(fmt.Sprintf("client_id: %s does not have a match for client_secret: %s", req.ClientID, req.ClientSecret))
		fmt.Println(err.Error())
		return
	}

	foundGrantType := false
	for _, element := range client.AllowedGrantTypes {

		foundGrantType = (req.GrantType == element)
		if foundGrantType {
			break
		}
	}
	if !foundGrantType {
		err = errors.New(fmt.Sprintf("client_id: %s is not authorized for grant_type: %s", req.ClientID, req.GrantType))
		fmt.Println(err.Error())
		return
	}

	return
}

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
	case "arbitrary_no_subject":
	case "arbitrary_identity":

	}
	resp := renderings.HealthCheckResponse{
		Status: "Should Never See this",
	}
	return c.JSON(http.StatusOK, resp)
}
