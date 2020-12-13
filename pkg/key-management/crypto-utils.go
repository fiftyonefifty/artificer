package keymanagement

import (
	"crypto"
	b64 "encoding/base64"
	"errors"

	"github.com/pascaldekloe/jwt"
)

var encoding = b64.RawURLEncoding
var errHashLink = errors.New("jwt: hash function not linked into binary")

func hashLookup(alg string, algs map[string]crypto.Hash) (crypto.Hash, error) {
	// availability check
	hash, ok := algs[alg]
	if !ok {
		return 0, jwt.AlgError(alg)
	}
	if !hash.Available() {
		return 0, errHashLink
	}
	return hash, nil
}

func HashByteArray(alg string, value *[]byte) (h *[]byte, err error) {
	hash, err := hashLookup(alg, jwt.ECDSAAlgs)
	if err != nil {
		return
	}

	digest := hash.New()
	digest.Write(*value)
	v := digest.Sum(nil)
	h = &v
	return
}
