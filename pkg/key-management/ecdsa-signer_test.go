package keymanagment

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"io"
	"testing"

	"github.com/pascaldekloe/jwt"
)

/*
private:
-----BEGIN PRIVATE KEY-----
MHcCAQEEIDOp7+7/aKFRBZOtKSXtE6sG71ALXqUaE8QcU/HzK1m7oAoGCCqGSM49
AwEHoUQDQgAEs0G8jc+uxcSnc7Ppfz4pkeKJ10sswOqKjexqUiB7UgbVOgMr1D59
bjEjJ8fQVghLZbrDHtpD7Xv9jD0qfcYvLw==
-----END PRIVATE KEY-----

public:
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEs0G8jc+uxcSnc7Ppfz4pkeKJ10ss
wOqKjexqUiB7UgbVOgMr1D59bjEjJ8fQVghLZbrDHtpD7Xv9jD0qfcYvLw==
-----END PUBLIC KEY-----
*/

type zr struct {
	io.Reader
}

// Read replaces the contents of dst with zeros.
func (z *zr) Read(dst []byte) (n int, err error) {
	for i := range dst {
		dst[i] = 0
	}
	return len(dst), nil
}

var zeroReader = &zr{}

func defaultTestTokenNoSig() (token []byte, err error) {
	var claims jwt.Claims
	claims.KeyID = "1324"
	claims.Set = make(map[string]interface{})
	claims.Set["A"] = "B"
	token, err = claims.FormatWithoutSign(jwt.ES256)
	return
}
func TestECDSASigner(t *testing.T) {
	ctx := context.Background()
	pemPrivate, pemPublic, err := PemGenerateECSDAP256()
	pemPrivate = privateEcdsaP256
	pemPublic = publicEcdsaP256
	privKey, pubKey := DecodeECDSA(pemPrivate, pemPublic)

	fmt.Printf("private:\n%s\npublic:\n%s\n", pemPrivate, pemPublic)
	signer := ECDSASigner{
		Key:       privKey,
		Algorithm: jwt.ES256,
		Kid:       "1324",
	}

	token, err := defaultTestTokenNoSig()
	if err != nil {
		t.Fail()
	}
	sToken := string(token)
	fmt.Printf("token:%s\n", sToken)
	h, err := HashByteArray(jwt.ES256, &token)
	sEnc := b64.StdEncoding.EncodeToString(*h)

	ksp := KeySignParameters{
		Algorithm: jwt.ES256,
		Value:     &sEnc,
	}

	result, err := signer.Sign(ctx, ksp)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("kid:%s, sig:%s\n", *result.Kid, *result.Result)

	token = append(token, '.')
	sig, err := b64.StdEncoding.DecodeString(*result.Result)

	token = append(token, sig...)

	sToken = string(token)
	fmt.Printf("token: %s\n", sToken)
	fmt.Println("token: eyJhbGciOiJFUzI1NiIsImtpZCI6IjEzMjQifQ.eyJBIjoiQiJ9.GuyWaYkwe_KpC3Mh8X76C0bTn8r4dtDyk-ORMO9LS2vyPz2gYApS80I-e7shhOI_qAsuybYYspGT_VTjEyiWqA")
	got, err := jwt.ECDSACheck(token, pubKey)
	fmt.Println(got.Set["A"])
	fmt.Println(got)

}
