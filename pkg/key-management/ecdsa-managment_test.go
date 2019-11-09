package keymanagment

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/lestrrat-go/jwx/jwk"
)

func TestA(t *testing.T) {
	pemPrivate, pemPublic, err := PemGenerateECSDAP256()
	if err != nil {
		t.Fail()
	}
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
