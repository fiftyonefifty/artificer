package handlers

import (
	"artificer/pkg/client/models"
	"artificer/pkg/keyvault"
	"artificer/pkg/util"
	"encoding/json"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
)

func buildArbitraryResourceOwnerClaims(req *ArbitraryResourceOwnerRequest) (err error, tokenBuildRequest keyvault.TokenBuildRequest) {

	if err = validateArbitraryResourceOwnerRequest(req); err != nil {
		return
	}

	client := req.Client

	var objmap map[string]interface{}
	err = json.Unmarshal([]byte(req.ArbitraryClaims), &objmap)

	// find anything that looks like our reservied models.NAMESPACE_NAME
	// remove it from the map
	reservedNames := []string{"scope", "aud", "amr", "nbf", "exp", "iss", "client_id", "sub", "auth_time", "idp", models.NAMESPACE_NAME}
	filterOutKeys := []string{}
	for k, _ := range objmap {
		if util.Contains(&reservedNames, k, true) {
			filterOutKeys = append(filterOutKeys, k)
		}
	}
	for _, k := range filterOutKeys {
		delete(objmap, k)
	}

	tokenBuildRequest.AccessTokenLifetime = client.AccessTokenLifetime
	if req.AccessTokenLifetime > 0 && req.AccessTokenLifetime < client.AccessTokenLifetime {
		tokenBuildRequest.AccessTokenLifetime = req.AccessTokenLifetime
	}
	utcNow := time.Now().UTC().Truncate(time.Minute)
	utcExpires := utcNow.Add(time.Second * time.Duration(tokenBuildRequest.AccessTokenLifetime))

	tokenBuildRequest.UtcExpires = &utcExpires
	tokenBuildRequest.UtcNotBefore = &utcNow

	tokenBuildRequest.Claims = jwt.Claims{
		// cover all registered fields
		Registered: jwt.Registered{},
		Set:        objmap,
	}
	if len(req.ArbitraryAudiences) > 0 {
		var arbAud []string

		err = json.Unmarshal([]byte(req.ArbitraryAudiences), &arbAud)
		if err != nil {
			return
		}
		tokenBuildRequest.Claims.Audiences = arbAud
	}

	for key, element := range tokenBuildRequest.Claims.Set {
		sArr := util.InterfaceArrayToStringArray(element)
		if sArr != nil {
			delete(tokenBuildRequest.Claims.Set, key)
			tokenBuildRequest.Claims.Set[key] = sArr
		}
	}

	tokenBuildRequest.Claims.Set["client_id"] = client.ClientID

	if client.AlwaysSendClientClaims {
		for _, element := range client.Claims {

			if tokenBuildRequest.Claims.Set[element.Type] == nil {
				tokenBuildRequest.Claims.Set[element.Type] = []string{}
			}
			set := tokenBuildRequest.Claims.Set[element.Type].([]string)
			set = append(set, element.Value)
			tokenBuildRequest.Claims.Set[element.Type] = set
		}
	}

	// TODO: Implement
	//offlineAccess:= false
	// deal with scopes
	if req.Scopes != nil && len(req.Subject) > 0 && req.Scopes["offline_access"] != nil {
		// we have a refresh_token request here.  TODO
		delete(req.Scopes, "offline_access")
		//		offlineAccess = true
	}
	if req.Scopes != nil && len(req.Scopes) > 0 {
		scopes := []string{}
		for scope, _ := range req.Scopes {
			scopes = append(scopes, scope)
		}
		tokenBuildRequest.Claims.Set["scope"] = scopes
	}

	for key, element := range tokenBuildRequest.Claims.Set {
		switch element.(type) {
		case []string:
			saElement := element.([]string)
			if len(saElement) == 1 {
				tokenBuildRequest.Claims.Set[key] = saElement[0]
			}
			break
		}
	}
	if len(req.ArbitraryAmrs) > 0 {
		var arbAmrs []string

		err = json.Unmarshal([]byte(req.ArbitraryAmrs), &arbAmrs)
		if err != nil {
			return
		}
		tokenBuildRequest.Claims.Set["amr"] = arbAmrs
	}

	if len(req.Subject) > 0 {
		tokenBuildRequest.Claims.Subject = req.Subject
	}

	return
}

func handleArbitraryResourceOwnerFlow(c echo.Context) (err error) {
	req := &ArbitraryResourceOwnerRequest{}
	if err = c.Bind(req); err != nil {
		return
	}
	req.ClientID = "<purposefully set to bad, use req.Client>"
	req.Client = c.Get("_client").(models.Client)

	err, tokenBuildRequest := buildArbitraryResourceOwnerClaims(req)
	if err != nil {
		return err
	}

	token, err := keyvault.MintToken(c, &tokenBuildRequest)
	if err != nil {
		return err
	}

	resp := ClientCredentialsResponse{
		AccessToken: token,
		ExpiresIn:   tokenBuildRequest.AccessTokenLifetime,
		TokenType:   "Bearer",
	}
	return c.JSON(http.StatusOK, resp)
}
