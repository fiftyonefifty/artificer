package keyvault

import (
	keymanagment "artificer/pkg/key-management"

	"context"

	azKeyvault "github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
)

type AzureKeyVaultSigner struct {
	BaseClient   azKeyvault.BaseClient
	VaultBaseURL string
	KeyName      string
}

func (signer AzureKeyVaultSigner) Sign(ctx context.Context, ksp keymanagment.KeySignParameters) (result keymanagment.KeyOperationResult, err error) {
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
	result = keymanagment.KeyOperationResult{
		Kid:    keyOperationResult.Kid,
		Result: keyOperationResult.Result,
	}
	return
}
