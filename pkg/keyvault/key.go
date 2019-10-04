package keyvault

import (
	"artificer/pkg/config"
	"artificer/pkg/iam"
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest/to"
)

func getKeysClient() keyvault.BaseClient {
	keyClient := keyvault.New()
	a, _ := iam.GetKeyvaultAuthorizer()
	keyClient.Authorizer = a
	keyClient.AddToUserAgent(config.UserAgent())
	return keyClient
}

func GetKeysVersion(ctx context.Context) (result keyvault.KeyListResultPage, err error) {
	keyClient := getKeysClient()
	var maxResults int32 = 10
	result, err = keyClient.GetKeyVersions(ctx, "https://P7KeyValut.vault.azure.net/", "P7IdentityServer4SelfSigned", &maxResults)
	return
}
func GetActiveKeysVersion(ctx context.Context) (finalResult []keyvault.KeyItem, err error) {
	keyClient := getKeysClient()
	var maxResults int32 = 10
	pageResult, er := keyClient.GetKeyVersions(ctx,
		"https://P7KeyValut.vault.azure.net/",
		"P7IdentityServer4SelfSigned",
		&maxResults)
	if er != nil {
		err = err
		return
	}

	utcNow := time.Now().UTC()

	for _, element := range pageResult.Values() {
		// element is the element from someSlice for where we are
		if *element.Attributes.Enabled {
			var keyExpire time.Time
			keyExpire = time.Time(*element.Attributes.Expires)

			if keyExpire.After(utcNow) {
				finalResult = append(finalResult, element)
			}
		}
	}

	for pageResult.NotDone() {
		er = pageResult.Next()
		if er != nil {
			err = err
			return
		}
		for _, element := range pageResult.Values() {
			// element is the element from someSlice for where we are
			if *element.Attributes.Enabled {
				var keyExpire time.Time
				keyExpire = time.Time(*element.Attributes.Expires)
				if keyExpire.After(utcNow) {
					finalResult = append(finalResult, element)
				}
			}
		}
	}
	return
}

// CreateKeyBundle creates a key in the specified keyvault
func CreateKey(ctx context.Context, vaultName, keyName string) (key keyvault.KeyBundle, err error) {
	vaultsClient := getVaultsClient()
	vault, err := vaultsClient.Get(ctx, config.GroupName(), vaultName)
	if err != nil {
		return
	}
	vaultURL := *vault.Properties.VaultURI

	keyClient := getKeysClient()
	return keyClient.CreateKey(
		ctx,
		vaultURL,
		keyName,
		keyvault.KeyCreateParameters{
			KeyAttributes: &keyvault.KeyAttributes{
				Enabled: to.BoolPtr(true),
			},
			KeySize: to.Int32Ptr(2048), // As of writing this sample, 2048 is the only supported KeySize.
			KeyOps: &[]keyvault.JSONWebKeyOperation{
				keyvault.Encrypt,
				keyvault.Decrypt,
			},
			Kty: keyvault.RSA,
		})
}
