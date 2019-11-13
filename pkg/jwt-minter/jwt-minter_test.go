package jwtMinter

import (
	keymanagement "artificer/pkg/key-management"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pascaldekloe/jwt"
)

func TestECSDA_JwtMinter(t *testing.T) {
	ctx := context.Background()
	pemPrivate, pemPublic, err := keymanagement.PemGenerateECSDAP256()
	if err != nil {
		t.Fail()
	}
	privKey, pubKey := keymanagement.DecodeECDSA(pemPrivate, pemPublic)

	jwtCtx := JwtMinterContext{
		Signer: keymanagement.ECDSASigner{
			Key: privKey,
			Kid: "1324",
		},
		Algorithm: jwt.ES256,
	}
	claims := keymanagement.DefaultTestClaims()
	notBefore := time.Now().UTC()
	expires := notBefore.Add(time.Second * 300)
	tbc := TokenBuildRequest{
		UtcNotBefore:        &notBefore,
		UtcExpires:          &expires,
		OfflineAccess:       false,
		AccessTokenLifetime: int(expires.Sub(notBefore).Seconds()),
		Claims:              claims,
	}
	token, err := jwtCtx.Create(ctx, &tbc)
	if err != nil {
		t.Fail()
	}
	fmt.Println(token)

	got, err := jwt.ECDSACheck([]byte(token), pubKey)
	fmt.Println(got.Set["A"])

}
