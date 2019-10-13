package util

import (
	echo "github.com/labstack/echo/v4"

	"artificer/pkg/keyvault"

	"time"

	"context"

	azKeyvault "github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/pascaldekloe/jwt"
	"github.com/spf13/viper"
)

func MintToken(c echo.Context, claims jwt.Claims, notBefore *time.Time, expires *time.Time) (token string, err error) {
	cachedItem, err := keyvault.GetCachedKeyVersions()
	if err != nil {
		return
	}

	baseUrl := GetBaseUrl(c)

	utcNow := time.Now().UTC()
	if notBefore == nil {
		notBefore = &utcNow
	}
	claims.Issued = jwt.NewNumericTime(utcNow)
	claims.NotBefore = jwt.NewNumericTime(*notBefore)
	if expires == nil {
		claims.Expires = nil
	} else {
		claims.Expires = jwt.NewNumericTime(*expires)
	}
	claims.KeyID = cachedItem.CurrentVersionId
	claims.Issuer = baseUrl

	keyClient, err := keyvault.GetKeyClient()
	if err != nil {
		return
	}
	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl")     //"https://P7KeyValut.vault.azure.net/"
	keyIdentifier := viper.GetString("keyVault.KeyIdentifier") //"P7IdentityServer4SelfSigned"
	ctx := context.Background()

	byteToken, err := keyClient.Sign2(ctx, &claims, azKeyvault.RS256, keyVaultUrl, keyIdentifier, claims.KeyID)
	if err != nil {
		return
	}
	token = string(byteToken)

	return
}
