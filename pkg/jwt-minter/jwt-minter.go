package jwtMinter

import (
	"time"

	keymanagement "artificer/pkg/key-management"
	"context"
	b64 "encoding/base64"

	"github.com/pascaldekloe/jwt"
)

type TokenBuildRequest struct {
	UtcNotBefore        *time.Time
	UtcExpires          *time.Time
	OfflineAccess       bool
	AccessTokenLifetime int
	Claims              jwt.Claims
}

type JwtMinterContext struct {
	Signer    keymanagement.Signer
	Algorithm string
}

type JwtMinter interface {
	Create(ctx context.Context, tokenBuildRequest *TokenBuildRequest) (token string, err error)
}

func (jwtCtx JwtMinterContext) Create(ctx context.Context, tokenBuildRequest *TokenBuildRequest) (token string, err error) {
	tokenWithoutSignature, err := tokenBuildRequest.Claims.FormatWithoutSign(jwtCtx.Algorithm)
	if err != nil {
		return
	}

	h, err := keymanagement.HashByteArray(jwt.ES256, &tokenWithoutSignature)
	sEnc := b64.StdEncoding.EncodeToString(*h)

	ksp := keymanagement.KeySignParameters{
		Algorithm: jwtCtx.Algorithm,
		Value:     &sEnc,
	}

	result, err := jwtCtx.Signer.Sign(ctx, ksp)
	if err != nil {
		return
	}

	tokenWithoutSignature = append(tokenWithoutSignature, '.')
	sig, err := b64.StdEncoding.DecodeString(*result.Result)

	tokenWithoutSignature = append(tokenWithoutSignature, sig...)
	token = string(tokenWithoutSignature)
	return
}
