package keyvault

import (
	"artificer/pkg/api/renderings"
	"artificer/pkg/config"
	"artificer/pkg/iam"
	"artificer/pkg/util"
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	azKeyvault "github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest/to"
	gocache "github.com/ghstahl/go-syncmap-cache"
	echo "github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/spf13/viper"
)

var (
	cache    = gocache.New(24*time.Hour, time.Hour)
	cacheKey = "85b75fb0-f120-4bfb-a0fe-f017cc72e41f"
)

type TokenBuildRequest struct {
	UtcNotBefore        *time.Time
	UtcExpires          *time.Time
	OfflineAccess       bool
	AccessTokenLifetime int
	Claims              jwt.Claims
}
type KeySignParameters struct {
	// Algorithm - The signing/verification algorithm identifier. For more information on possible algorithm types, see JsonWebKeySignatureAlgorithm. Possible values include: 'PS256', 'PS384', 'PS512', 'RS256', 'RS384', 'RS512', 'RSNULL', 'ES256', 'ES384', 'ES512', 'ES256K'
	Algorithm string `json:"alg,omitempty"`
	// Value - a URL-encoded base64 string
	Value *string `json:"value,omitempty"`
}
type KeyOperationResult struct {
	// Kid - READ-ONLY; Key identifier
	Kid *string `json:"kid,omitempty"`
	// Result - READ-ONLY; a URL-encoded base64 string
	Result *string `json:"value,omitempty"`
}
type Signer interface {
	Sign(ctx context.Context, ksp KeySignParameters) (result KeyOperationResult, err error)
}
type AzureKeyVaultSigner struct {
	BaseClient   azKeyvault.BaseClient
	VaultBaseURL string
	KeyName      string
}

func (signer AzureKeyVaultSigner) Sign(ctx context.Context, ksp KeySignParameters) (result KeyOperationResult, err error) {
	var azAlg azKeyvault.JSONWebKeySignatureAlgorithm = azKeyvault.JSONWebKeySignatureAlgorithm(ksp.Algorithm)

	cachedItem, err := GetCachedKeyVersions()
	if err != nil {
		return
	}
	keyOperationResult, err := signer.BaseClient.Sign(ctx, signer.VaultBaseURL, signer.KeyName, cachedItem.CurrentVersionId,
		azKeyvault.KeySignParameters{
			Algorithm: azAlg,
			Value:     ksp.Value,
		})
	if err != nil {
		return
	}
	result = KeyOperationResult{
		Kid:    keyOperationResult.Kid,
		Result: keyOperationResult.Result,
	}
	return
}

type BaseClient2 struct {
	azKeyvault.BaseClient
	Signer Signer
}

func newAzureKeyVaultBaseClient2(base azKeyvault.BaseClient) BaseClient2 {
	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl") //"https://P7KeyValut.vault.azure.net/"
	keyName := viper.GetString("keyVault.KeyIdentifier")   //"P7IdentityServer4SelfSigned"

	return BaseClient2{
		BaseClient: base,
		Signer: AzureKeyVaultSigner{
			BaseClient:   base,
			VaultBaseURL: keyVaultUrl,
			KeyName:      keyName,
		},
	}
}

func getKeysClient() azKeyvault.BaseClient {
	keyClient := azKeyvault.New()
	a, _ := iam.GetKeyvaultAuthorizer()
	keyClient.Authorizer = a
	keyClient.AddToUserAgent(config.UserAgent())
	return keyClient
}

func GetSecret(name string) (result keyvault.SecretBundle, err error) {
	ctx := context.Background()
	keyClient := getKeysClient()
	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl") //"https://P7KeyValut.vault.azure.net/"

	return keyClient.GetSecret(ctx, keyVaultUrl, name, "")
}

func GetKeysVersion(ctx context.Context) (result azKeyvault.KeyListResultPage, err error) {

	keyClient := getKeysClient()
	var maxResults int32 = 10
	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl")     //"https://P7KeyValut.vault.azure.net/"
	keyIdentifier := viper.GetString("keyVault.KeyIdentifier") //"P7IdentityServer4SelfSigned"

	result, err = keyClient.GetKeyVersions(ctx, keyVaultUrl, keyIdentifier, &maxResults)
	return
}

func MintToken(c echo.Context, tokenBuildRequest *TokenBuildRequest) (token string, err error) {
	cachedItem, err := GetCachedKeyVersions()
	if err != nil {
		return
	}

	baseUrl := util.GetBaseUrl(c)

	utcNow := time.Now().UTC().Truncate(time.Second)
	if tokenBuildRequest.UtcNotBefore == nil {
		tokenBuildRequest.UtcNotBefore = &utcNow
	}
	tokenBuildRequest.Claims.Issued = jwt.NewNumericTime(utcNow)
	tokenBuildRequest.Claims.NotBefore = jwt.NewNumericTime(*tokenBuildRequest.UtcNotBefore)
	if tokenBuildRequest.UtcExpires == nil {
		tokenBuildRequest.Claims.Expires = nil
	} else {
		tokenBuildRequest.Claims.Expires = jwt.NewNumericTime(*tokenBuildRequest.UtcExpires)
	}
	tokenBuildRequest.Claims.KeyID = cachedItem.CurrentVersionId
	tokenBuildRequest.Claims.Issuer = baseUrl

	if tokenBuildRequest.Claims.Audiences == nil {
		tokenBuildRequest.Claims.Audiences = []string{}
	}
	tokenBuildRequest.Claims.Audiences = append(tokenBuildRequest.Claims.Audiences, tokenBuildRequest.Claims.Issuer)

	keyClient, err := GetKeyClient()
	if err != nil {
		return
	}
	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl")     //"https://P7KeyValut.vault.azure.net/"
	keyIdentifier := viper.GetString("keyVault.KeyIdentifier") //"P7IdentityServer4SelfSigned"
	ctx := context.Background()

	byteToken, err := keyClient.Sign2(ctx, &tokenBuildRequest.Claims, azKeyvault.RS256, keyVaultUrl, keyIdentifier, tokenBuildRequest.Claims.KeyID)
	if err != nil {
		return
	}
	token = string(byteToken)

	return
}
func RSA256AzureSign(ctx context.Context, data []byte) (kid *string, signature *string, err error) {
	keyClient := getKeysClient()

	sEnc := util.ByteArraySha256Encode64(data)

	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl")     //"https://P7KeyValut.vault.azure.net/"
	keyIdentifier := viper.GetString("keyVault.KeyIdentifier") //"P7IdentityServer4SelfSigned"

	keyOperationResult, err := keyClient.Sign(ctx, keyVaultUrl, keyIdentifier, "", keyvault.KeySignParameters{
		Algorithm: azKeyvault.RS256,
		Value:     &sEnc,
	})
	if err != nil {
		return
	}
	return keyOperationResult.Kid, keyOperationResult.Result, nil
}

func (keyClient *BaseClient2) Sign2(
	ctx context.Context,
	claims *jwt.Claims,
	alg azKeyvault.JSONWebKeySignatureAlgorithm,
	vaultBaseURL string,
	keyName string,
	keyVersion string) (token []byte, err error) {

	tokenWithoutSignature, err := claims.FormatWithoutSign(string(alg))
	if err != nil {
		return nil, err
	}
	sEnc := util.ByteArraySha256Encode64(tokenWithoutSignature)

	/*
		keyOperationResult, err := keyClient.Sign(ctx, vaultBaseURL, keyName, keyVersion, azKeyvault.KeySignParameters{
			Algorithm: alg,
			Value:     &sEnc,
		})
	*/
	keyOperationResult, err := keyClient.Signer.Sign(ctx, KeySignParameters{
		Algorithm: string(alg),
		Value:     &sEnc,
	})
	if err != nil {
		return
	}

	token = append(tokenWithoutSignature, '.')
	token = append(token, []byte(*keyOperationResult.Result)...)
	return token, nil
}

func (keyClient *BaseClient2) GetActiveKeysVersion2(ctx context.Context, keyVaultUrl string, keyIdentifier string) (finalResult []azKeyvault.KeyBundle, currentKeyBundle azKeyvault.KeyBundle, err error) {

	// Length requirements defined by 2.2.2.9.1 RSA Private Key BLOB (https://msdn.microsoft.com/en-us/library/cc250013.aspx).
	/*
		PubExp (4 bytes): Length MUST be 4 bytes.

		This field MUST be present as an unsigned integer in little-endian format.

		The value of this field MUST be the RSA public key exponent for this key. The client SHOULD set this value to 65,537.

		E is comming back as an Base64Url encoded byte[] of size 3.
	*/

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

func GetKeyClient() (keyClient BaseClient2, err error) {

	baseClient := getKeysClient()
	keyClient = newAzureKeyVaultBaseClient2(baseClient)
	err = nil
	return
}

func GetActiveKeysVersion(ctx context.Context) (finalResult []azKeyvault.KeyBundle, currentKeyBundle azKeyvault.KeyBundle, err error) {

	keyVaultUrl := viper.GetString("keyVault.KeyVaultUrl")     //"https://P7KeyValut.vault.azure.net/"
	keyIdentifier := viper.GetString("keyVault.KeyIdentifier") //"P7IdentityServer4SelfSigned"
	//keyClient := getKeysClient()

	baseClient2, err := GetKeyClient()
	if err != nil {
		return
	}
	finalResult, currentKeyBundle, err = baseClient2.GetActiveKeysVersion2(ctx, keyVaultUrl, keyIdentifier)

	return
}

// CreateKeyBundle creates a key in the specified keyvault
func CreateKey(ctx context.Context, vaultName, keyName string) (key azKeyvault.KeyBundle, err error) {
	vaultsClient := getVaultsClient()
	vault, err := vaultsClient.Get(ctx, config.BaseGroupName(), vaultName)
	if err != nil {
		return
	}
	vaultURL := *vault.Properties.VaultURI

	keyClient := getKeysClient()
	return keyClient.CreateKey(
		ctx,
		vaultURL,
		keyName,
		azKeyvault.KeyCreateParameters{
			KeyAttributes: &azKeyvault.KeyAttributes{
				Enabled: to.BoolPtr(true),
			},
			KeySize: to.Int32Ptr(2048), // As of writing this sample, 2048 is the only supported KeySize.
			KeyOps: &[]azKeyvault.JSONWebKeyOperation{
				azKeyvault.Encrypt,
				azKeyvault.Decrypt,
			},
			Kty: azKeyvault.EC,
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

type CachedKeyVersions struct {
	CurrentKeyBundle                         azKeyvault.KeyBundle
	CurrentVersionId                         string
	WellKnownOpenidConfigurationJwksResponse renderings.WellKnownOpenidConfigurationJwksResponse
}

func GetCachedKeyVersions() (cachedResponse CachedKeyVersions, err error) {

	var cachedItem interface{}
	var found bool

	cachedItem, found = cache.Get(cacheKey)

	if !found {
		err = DoKeyVaultBackground()
		if err != nil {
			log.Fatalf("failed to DoKeyvaultBackground: %v\n", err.Error())
			return
		}
		cachedItem, found = cache.Get(cacheKey)
		if !found {
			err = errors.New("critical failure to DoKeyvaultBackground")
			log.Fatalln(err.Error())
			return
		}
	}
	cachedResponse = cachedItem.(CachedKeyVersions)
	return
}

func DoKeyVaultBackground() (err error) {
	now := time.Now().UTC()
	fmt.Println(fmt.Sprintf("Start-DoKeyvaultBackground:%s", now))
	ctx := context.Background()
	activeKeys, currentKeyBundle, err := GetActiveKeysVersion(ctx)
	if err != nil {
		return
	}
	resp := renderings.WellKnownOpenidConfigurationJwksResponse{}

	for _, element := range activeKeys {

		parts := strings.Split(*element.Key.Kid, "/")
		lastItemVersion := parts[len(parts)-1]

		jwk := renderings.JwkResponse{}
		jwk.Kid = lastItemVersion
		jwk.Kty = string(element.Key.Kty)
		jwk.N = *element.Key.N
		jwk.E = *element.Key.E
		jwk.Alg = "RSA256"
		jwk.Use = "sig"
		resp.Keys = append(resp.Keys, jwk)
	}

	parts := strings.Split(*currentKeyBundle.Key.Kid, "/")
	lastItemVersion := parts[len(parts)-1]

	cacheItem := CachedKeyVersions{
		WellKnownOpenidConfigurationJwksResponse: resp,
		CurrentKeyBundle:                         currentKeyBundle,
		CurrentVersionId:                         lastItemVersion,
	}

	cache.Set(cacheKey, cacheItem, gocache.NoExpiration)
	fmt.Println(fmt.Sprintf("Success-DoKeyvaultBackground:%s", now))
	return
}
