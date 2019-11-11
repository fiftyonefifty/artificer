package keymanagment

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/lestrrat-go/jwx/jwk"
)

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
