package handlers

import (
	"artificer/pkg/api/models"
	"artificer/pkg/config"
	"artificer/pkg/keyvault"
	"artificer/pkg/util"
	"encoding/json"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
)

func buildArbitraryNoSubjectClaims(req *ArbitraryNoSubjectRequest) (err error, utcNotBefore time.Time, utcExpires time.Time, accessTokenLifetime int, claims jwt.Claims) {

	if err = validateArbitraryNoSubjectRequest(req); err != nil {
		return
	}

	var client *models.Client
	client = config.ClientMap[req.ClientID]

	var objmap map[string]interface{}
	err = json.Unmarshal([]byte(req.ArbitraryClaims), &objmap)

	// find anything that looks like our reservied models.NAMESPACE_NAME
	// remove it from the map
	reservedNames := []string{"aud", "amr", "nbf", "exp", "iss", "client_id", "sub", "auth_time", "idp", models.NAMESPACE_NAME}
	filterOutKeys := []string{}
	for k, _ := range objmap {
		if util.Contains(&reservedNames, k, true) {
			filterOutKeys = append(filterOutKeys, k)
		}
	}
	for _, k := range filterOutKeys {
		delete(objmap, k)
	}

	accessTokenLifetime = client.AccessTokenLifetime
	if req.AccessTokenLifetime > 0 && req.AccessTokenLifetime < client.AccessTokenLifetime {
		accessTokenLifetime = req.AccessTokenLifetime
	}
	utcNow := time.Now().UTC().Truncate(time.Minute)
	utcExpires = utcNow.Add(time.Second * time.Duration(accessTokenLifetime))
	utcNotBefore = utcNow

	claims = jwt.Claims{
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
		claims.Audiences = arbAud
	}

	for key, element := range claims.Set {
		er, sArr := util.InterfaceArrayToStringArray(element)
		if er != nil {
			err = er
			return
		}
		delete(claims.Set, key)
		claims.Set[key] = sArr
	}

	claims.Set["client_id"] = client.ClientID

	if client.AlwaysSendClientClaims {
		for _, element := range client.Claims {

			if claims.Set[element.Type] == nil {
				claims.Set[element.Type] = []string{}
			}
			set := claims.Set[element.Type].([]string)
			set = append(set, element.Value)
			claims.Set[element.Type] = set
		}
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
	if len(req.ArbitraryAmrs) > 0 {
		var arbAmrs []string

		err = json.Unmarshal([]byte(req.ArbitraryAmrs), &arbAmrs)
		if err != nil {
			return
		}
		claims.Set["amr"] = arbAmrs
	}
	if len(req.CustomPayload) > 0 {
		var objmap interface{}
		err = json.Unmarshal([]byte(req.CustomPayload), &objmap)
		if err != nil {
			return
		}
		claims.Set["custom_payload"] = objmap
	}

	return
}

func handleArbitraryNoSubjectFlow(c echo.Context) (err error) {
	req := &ArbitraryNoSubjectRequest{}
	if err = c.Bind(req); err != nil {
		return
	}
	err, utcNotBefore, utcExpires, accessTokenLifetime, claims := buildArbitraryNoSubjectClaims(req)
	if err != nil {
		return err
	}

	token, err := keyvault.MintToken(c, claims, &utcNotBefore, &utcExpires)
	if err != nil {
		return err
	}

	resp := ClientCredentialsResponse{
		AccessToken: token,
		ExpiresIn:   accessTokenLifetime,
		TokenType:   "Bearer",
	}
	return c.JSON(http.StatusOK, resp)
}
