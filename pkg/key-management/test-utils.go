package keymanagement

import (
	"io"

	"github.com/pascaldekloe/jwt"
)

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

func DefaultTestClaims() (claims jwt.Claims) {

	claims.KeyID = "1324"
	claims.Set = make(map[string]interface{})
	claims.Set["A"] = "B"
	return
}

func DefaultTestTokenNoSig() (token []byte, err error) {
	var claims jwt.Claims = DefaultTestClaims()
	token, err = claims.FormatWithoutSign(jwt.ES256)
	return
}
