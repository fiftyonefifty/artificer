package handlers

import (
	"artificer/pkg/keyvault"
	"artificer/pkg/util"
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
	utcNow := time.Now().UTC()
	utcExpires := utcNow.Add(time.Minute * 31)
	utcNotBefore := utcNow.Add(time.Minute * 1)

	claims := jwt.Claims{
		// cover all registered fields
		Registered: jwt.Registered{
			Subject:   "b",
			Audiences: []string{"c"},
			ID:        "d",
		},
		Set: make(map[string]interface{}),
	}
	claims.Set["pirate"] = "jack"
	claims.Set["primes"] = []int{2, 3, 5, 7, 11, 13}
	claims.Set["roles"] = []string{"admin", "super-duper"}

	token, err := util.MintToken(c, claims, &utcNotBefore, &utcExpires)

	return c.JSON(http.StatusOK, token)
}
