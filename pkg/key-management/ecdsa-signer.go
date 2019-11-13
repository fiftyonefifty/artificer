package keymanagment

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	b64 "encoding/base64"
	"strconv"
)

type ECDSASigner struct {
	Key       *ecdsa.PrivateKey
	Algorithm string
	Kid       string
}

func (signer ECDSASigner) Sign(ctx context.Context, ksp KeySignParameters) (result KeyOperationResult, err error) {

	dEnc, err := b64.StdEncoding.DecodeString(*ksp.Value)
	if err != nil {
		return
	}
	ioReader := rand.Reader
	ioReader = zeroReader
	r, s, err := ecdsa.Sign(ioReader, signer.Key, dEnc)
	if err != nil {
		return
	}
	paramLen := (signer.Key.Curve.Params().BitSize + 7) / 8
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

	encSig := b64.StdEncoding.EncodeToString(sig)
	result = KeyOperationResult{
		Kid:    &signer.Kid,
		Result: &encSig,
	}

	return
}
