package handlers

import (
	"artificer/pkg/client/clientContext"
	"artificer/pkg/client/models"
	"artificer/pkg/keyvault"
	"artificer/pkg/util"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
)

func handleClientCredentialsFlow(ctx context.Context, c echo.Context) error {
	req := &ClientCredentialsRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}
	var client models.Client
	var ok bool
	client, ok = clientContext.FromContext(ctx)
	if !ok {
		return errors.New("context.Context doesn't contain client object")
	}
	//	client = c.Get("_client").(models.Client)

	//	_, client := clientStore.GetClient(req.ClientID)

	scope := strings.TrimSpace(req.Scope)

	var allowedScopes []string
	if len(scope) == 0 {
		// OK, since scope is optional
		// by spec if nothing is asked for all is returned.
		allowedScopes = client.AllowedScopes
		if allowedScopes == nil {
			allowedScopes = []string{}
		}
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

	claims.Set["client_id"] = client.ClientID
	claims.Set["scope"] = allowedScopes

	if client.AlwaysSendClientClaims {
		for _, element := range client.Claims {

			if claims.Set[element.Type] == nil {
				claims.Set[element.Type] = []string{}
			}
			claims.Set[element.Type] = append(claims.Set[element.Type].([]string), element.Value)
		}
	}
	if claims.Set["scope"] != nil {
		claims.Registered.Audiences = claims.Set["scope"].([]string)
	}

	for key, element := range claims.Set {
		switch element.(type) {
		case []string:
			saElement := element.([]string)
			if len(saElement) == 1 {
				claims.Set[key] = saElement[0]
			}
			break
		}
	}
	tokenBuildRequest := keyvault.TokenBuildRequest{
		Claims:       claims,
		UtcNotBefore: &utcNotBefore,
		UtcExpires:   &utcExpires,
	}
	tokenBuildRequest.Claims = claims
	token, err := keyvault.MintToken(c, &tokenBuildRequest)
	if err != nil {
		return err
	}

	resp := ClientCredentialsResponse{
		AccessToken: token,
		ExpiresIn:   client.AccessTokenLifetime,
		TokenType:   "Bearer",
	}
	return c.JSON(http.StatusOK, resp)
}
