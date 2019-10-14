package handlers

import (
	"artificer/pkg/api/models"
	"artificer/pkg/api/renderings"
	"artificer/pkg/config"
	"artificer/pkg/keyvault"
	"artificer/pkg/util"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
)

type (
	TokenRequest struct {
		ClientID     string `json:"client_id" form:"client_id" query:"client_id"`
		ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
		GrantType    string `json:"grant_type" form:"grant_type" query:"grant_type"`
	}
	ClientCredentialsRequest struct {
		TokenRequest
		Scope string `json:"scope" form:"scope" query:"scope"`
	}
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
func handleClientCredentialsFlow(c echo.Context) error {
	req := &ClientCredentialsRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}
	var client *models.Client
	client = config.ClientMap[req.ClientID]

	scope := strings.TrimSpace(req.Scope)

	var allowedScopes []string
	if len(scope) == 0 {
		// OK, since scope is optional
		// by spec if nothing is asked for all is returned.
		allowedScopes = client.AllowedScopes
	} else {
		// this one gets weird, if you ask for something and that something doesn't exit you get NOTHING
		// which is odd, because asking for nothing gets you everything as can be seen if scope request is empty
		splitScopes := strings.Split(scope, " ")
		allowedScopes = util.IntersectionStringArray(splitScopes, client.AllowedScopes)
	}

	utcNow := time.Now().UTC().Truncate(time.Minute)
	utcExpires := utcNow.Add(time.Second * time.Duration(client.AccessTokenLifetime))
	utcNotBefore := utcNow

	claims := jwt.Claims{
		// cover all registered fields
		Registered: jwt.Registered{
			Audiences: allowedScopes,
		},
		Set: make(map[string]interface{}),
	}
	claims.Set["scope"] = allowedScopes
	claims.Set["client_id"] = client.ClientID
	/*
		tmpMap := make(map[string][]string)
			if client.AlwaysSendClientClaims {
				for _, element := range client.Claims {
					if element.Type == "scope" {
						continue
					}
					if tmpMap[element.Type] == nil {
						tmpMap[element.Type] = []string{}
					}
					append(tmpMap[element.Type], element.Value)
				}
			}
	*/

	token, err := keyvault.MintToken(c, claims, &utcNotBefore, &utcExpires)
	if err != nil {
		return err
	}

	resp := renderings.HealthCheckResponse{
		Status: fmt.Sprintf("UP: %s", token),
	}
	return c.JSON(http.StatusOK, resp)
}
