package crypto

import (
	"errors"
	"fmt"
	"io"
)

type ICodec interface {
	Name() string
	Encrypt(buf []byte, key []byte) ([]byte, error)
	Decrypt(buf []byte, key []byte) ([]byte, error)
}

var codecMapping = map[string]ICodec{
	CodecAesCBC:              NewAesCBC(),
	CodecAesCFB:              NewAesCFB(),
	CodecAesCTR:              NewAesCTR(),
	CodecAesGCM:              NewAesGCM(),
	CodecAesOFB:              NewAesOFB(),
	CodecAesChacha20Poly1305: NewChacha20Poly1305(),
}

var (
	errCiphertextTooShort = errors.New("ciphertext is too short")
	errInvalidBlockLength = errors.New("ciphertext length is not a multiple of the block size")
	errInvalidPadding     = errors.New("invalid padding")
	errFrameTooLarge      = errors.New("encrypted frame is too large")
)

func wrapCryptoReadError(operation string, err error) error {
	if errors.Is(err, io.EOF) {
		return io.EOF
	}
	return fmt.Errorf("%s: %w", operation, err)
}

func FindCodec(name string) (ICodec, bool) {
	if c, ok := codecMapping[name]; ok {
		return c, true
	}
	return nil, false
}
