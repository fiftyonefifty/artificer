package keymanagement

import (
	"artificer/pkg/util"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pascaldekloe/jwt"
)

func TestFormatWithoutSign(t *testing.T) {

	token, err := defaultTestTokenNoSig()

	if err != nil {
		t.Fatal("sign error:", err)
	}

	sToken := string(token)
	fmt.Println(sToken)
	encodedToken := util.ByteArraySha256Encode64(token)
	fmt.Println(encodedToken)
	pemPrivate, pemPublic, err := PemGenerateECSDAP256()
	pemPrivate = privateEcdsaP256
	pemPublic = publicEcdsaP256

	priv2, pub2 := DecodeECDSA(pemPrivate, pemPublic)

	paramLen := (priv2.Curve.Params().BitSize + 7) / 8

	h, err := HashByteArray(jwt.ES256, &token)

	r, s, err := ecdsa.Sign(zeroReader, priv2, *h)

	sig := make([]byte, encoding.EncodedLen(paramLen*2))
	i := len(sig)
	for _, word := range s.Bits() {
		for bitCount := strconv.IntSize; bitCount > 0; bitCount -= 8 {
			i--
			sig[i] = byte(word)
			word >>= 8
		}
	}
	// i might have exceeded paramLen due to the word size
	i = len(sig) - paramLen
	for _, word := range r.Bits() {
		for bitCount := strconv.IntSize; bitCount > 0; bitCount -= 8 {
			i--
			sig[i] = byte(word)
			word >>= 8
		}
	}

	// encoder won't overhaul source space
	encoding.Encode(sig, sig[len(sig)-2*paramLen:])

	token = append(token, '.')
	token = append(token, sig...)
	sToken = string(token)
	fmt.Println(sToken)

	got, err := jwt.ECDSACheck(token, pub2)
	fmt.Println(got.Set["A"])
	fmt.Println(got)
}

func validatePem(t *testing.T, pemPrivate string, pemPublic string) {

	fmt.Println(pemPrivate)
	fmt.Println(pemPublic)

	pubKey := DecodeECDSAPublic(pemPublic)

	key, err := jwk.New(pubKey)
	if err != nil {
		log.Printf("failed to create JWK: %s", err)
		return
	}
	jsonbuf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		log.Printf("failed to generate JSON: %s", err)
		return
	}
	fmt.Println(string(jsonbuf))
}
func Test_PemGenerateECSDAP224(t *testing.T) {
	pemPrivate, pemPublic, err := PemGenerateECSDAP224()
	if err != nil {
		t.Fail()
	}
	validatePem(t, pemPrivate, pemPublic)
}
func Test_PemGenerateECSDAP256(t *testing.T) {
	pemPrivate, pemPublic, err := PemGenerateECSDAP256()
	if err != nil {
		t.Fail()
	}
	validatePem(t, pemPrivate, pemPublic)
}
func Test_PemGenerateECSDAP384(t *testing.T) {
	pemPrivate, pemPublic, err := PemGenerateECSDAP384()
	if err != nil {
		t.Fail()
	}
	validatePem(t, pemPrivate, pemPublic)
}
func Test_PemGenerateECSDAP521(t *testing.T) {
	pemPrivate, pemPublic, err := PemGenerateECSDAP521()
	if err != nil {
		t.Fail()
	}
	validatePem(t, pemPrivate, pemPublic)
}
