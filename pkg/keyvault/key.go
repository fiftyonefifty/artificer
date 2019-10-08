package keyvault

import (
	"artificer/pkg/config"
	"artificer/pkg/iam"
	"context"
	"crypto"
	b64 "encoding/base64"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/spf13/viper"
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

func RSA256AzureSign(ctx context.Context, data []byte) (kid *string, signature *string, err error) {
	keyClient := getKeysClient()
	digest := crypto.SHA256.New()
	digest.Write(data)
	h := digest.Sum(nil)
	sEnc := b64.StdEncoding.EncodeToString(h)

	keyOperationResult, err := keyClient.Sign(ctx, "https://P7KeyValut.vault.azure.net/", "P7IdentityServer4SelfSigned", "", keyvault.KeySignParameters{
		Algorithm: keyvault.RS256,
		Value:     &sEnc,
	})
	if err != nil {
		return
	}
	return keyOperationResult.Kid, keyOperationResult.Result, nil
}

func GetActiveKeysVersion(ctx context.Context) (finalResult []keyvault.KeyBundle, currentKeyBundle keyvault.KeyBundle, err error) {

	// Length requirements defined by 2.2.2.9.1 RSA Private Key BLOB (https://msdn.microsoft.com/en-us/library/cc250013.aspx).
	/*
		PubExp (4 bytes): Length MUST be 4 bytes.

		This field MUST be present as an unsigned integer in little-endian format.

		The value of this field MUST be the RSA public key exponent for this key. The client SHOULD set this value to 65,537.

		E is comming back as an Base64Url encoded byte[] of size 3.
	*/

	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl")     //"https://P7KeyValut.vault.azure.net/"
	keyIdentifier := viper.GetString("keyVault.KeyIdentifier") //"P7IdentityServer4SelfSigned"
	keyClient := getKeysClient()

	var maxResults int32 = 10
	pageResult, err := keyClient.GetKeyVersions(ctx,
		keyVaultUrl,
		keyIdentifier,
		&maxResults)
	if err != nil {
		return
	}

	utcNow := time.Now().UTC()
	for {
		for _, element := range pageResult.Values() {
			// element is the element from someSlice for where we are
			if *element.Attributes.Enabled {
				var keyExpire time.Time
				keyExpire = time.Time(*element.Attributes.Expires)
				if keyExpire.After(utcNow) {
					parts := strings.Split(*element.Kid, "/")
					lastItemVersion := parts[len(parts)-1]

					keyBundle, er := keyClient.GetKey(ctx,
						keyVaultUrl,
						keyIdentifier,
						lastItemVersion)
					if er != nil {
						err = er
						return
					}
					fixedE := fixE(*keyBundle.Key.E)
					*keyBundle.Key.E = fixedE
					finalResult = append(finalResult, keyBundle)
				}
			}
		}
		if !pageResult.NotDone() {
			break
		}
		err = pageResult.Next()
		if err != nil {
			return
		}
	}

	sort.Slice(finalResult[:], func(i, j int) bool {
		notBeforeA := time.Time(*finalResult[i].Attributes.NotBefore)
		notBeforeB := time.Time(*finalResult[j].Attributes.NotBefore)

		return notBeforeA.After(notBeforeB)
	})

	for _, element := range finalResult {
		notVBefore := time.Time(*element.Attributes.NotBefore)
		if notVBefore.Before(utcNow) {
			currentKeyBundle = element
			break
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

func fixE(base64EncodedE string) string {
	sDec, _ := b64.StdEncoding.DecodeString(base64EncodedE)
	sDec = forceByteArrayLength(sDec, 4)
	sEnc := b64.StdEncoding.EncodeToString(sDec)
	parts := strings.Split(sEnc, "=")
	sEnc = parts[0]
	return sEnc
}

func forceByteArrayLength(slice []byte, requireLength int) []byte {
	n := len(slice)
	if n >= requireLength {
		return slice
	}
	newSlice := make([]byte, requireLength)
	offset := requireLength - n
	copy(newSlice[offset:], slice)
	slice = newSlice
	return slice
}
