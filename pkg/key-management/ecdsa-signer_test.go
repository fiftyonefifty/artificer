package keymanagement

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"testing"

	"github.com/pascaldekloe/jwt"
)

func TestECDSASigner(t *testing.T) {
	ctx := context.Background()
	pemPrivate, pemPublic, err := PemGenerateECSDAP256()
	pemPrivate = privateEcdsaP256
	pemPublic = publicEcdsaP256
	privKey, pubKey := DecodeECDSA(pemPrivate, pemPublic)

	fmt.Printf("private:\n%s\npublic:\n%s\n", pemPrivate, pemPublic)
	signer := ECDSASigner{
		Key: privKey,
		Kid: "1324",
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
		Value: &sEnc,
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
