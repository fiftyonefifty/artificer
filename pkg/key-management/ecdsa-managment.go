package keymanagment

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func EncodeECDSA(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func DecodeECDSA(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey := DecodeECDSAPrivate(pemEncoded)
	publicKey := DecodeECDSAPublic(pemEncodedPub)
	return privateKey, publicKey
}
func DecodeECDSAPublic(pemEncodedPub string) *ecdsa.PublicKey {

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}
func DecodeECDSAPrivate(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}

func PemGenerateECSDAP256() (pemPrivate string, pemPublic string, err error) {

	pubkeyCurve := elliptic.P256()                                 //see http://golang.org/pkg/crypto/elliptic/#P256
	privateKey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
	if err != nil {
		return
	}
	var publicKey ecdsa.PublicKey
	publicKey = privateKey.PublicKey

	pemEncodedPrivate, pemEncodedPub := EncodeECDSA(privateKey, &publicKey)

	priv2, pub2 := DecodeECDSA(pemEncodedPrivate, pemEncodedPub)

	pemEncodedPrivate2, pemEncodedPub2 := EncodeECDSA(priv2, pub2)

	if pemEncodedPrivate != pemEncodedPrivate2 {
		fmt.Println("Private keys do not match.")
		err = errors.New("Private keys do not match.")
		return
	}
	if pemEncodedPub != pemEncodedPub2 {
		fmt.Println("Public keys do not match.")
		err = errors.New("Public keys do not match.")
		return
	}
	pemPrivate = pemEncodedPrivate
	pemPublic = pemEncodedPub
	return
}
func mustParseECKey(s string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		panic("invalid PEM")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}
