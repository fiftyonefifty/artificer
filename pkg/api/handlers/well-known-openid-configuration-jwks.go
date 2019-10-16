package handlers

import (
	"artificer/pkg/keyvault"
	"net/http"

	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
)

func WellKnownOpenidConfigurationJwks(c echo.Context) (err error) {
	cachedItem, err := keyvault.GetCachedKeyVersions()
	if err != nil {
		return
	}
	return c.JSON(http.StatusOK, cachedItem.WellKnownOpenidConfigurationJwksResponse)
}

func MintTestToken(c echo.Context) (err error) {
	utcNow := time.Now().UTC().Truncate(time.Minute)
	utcExpires := utcNow.Add(time.Minute * 31)
	utcNotBefore := utcNow.Add(time.Minute * 1)

	claims := jwt.Claims{
		// cover all registered fields
		Registered: jwt.Registered{
			//	Subject:   "b",
			Audiences: []string{"aud1", "aud2"},
			ID:        "d",
		},
		Set: make(map[string]interface{}),
	}
	claims.Set["pirate"] = "jack"
	claims.Set["primes"] = []int{2, 3, 5, 7, 11, 13}
	claims.Set["roles"] = []string{"admin", "super-duper"}
	claims.Set["scope"] = []string{"aud1", "aud2"}
	tokenBuildRequest := keyvault.TokenBuildRequest{
		Claims:       claims,
		UtcNotBefore: &utcNotBefore,
		UtcExpires:   &utcExpires,
	}

	token, err := keyvault.MintToken(c, &tokenBuildRequest)

	return c.JSON(http.StatusOK, token)
}
