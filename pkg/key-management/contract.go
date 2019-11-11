package keymanagment

import (
	"context"
)

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
