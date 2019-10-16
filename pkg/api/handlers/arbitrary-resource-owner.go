package handlers

import (
	"artificer/pkg/keyvault"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func handleArbitraryResourceOwnerFlow(c echo.Context) (err error) {
	req := &ArbitraryResourceOwnerRequest{}
	if err = c.Bind(req); err != nil {
		return
	}
	reqB := req.ArbitraryNoSubjectRequest

	err, utcNotBefore, utcExpires, accessTokenLifetime, claims := buildArbitraryNoSubjectClaims(&reqB)
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
