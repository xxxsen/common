package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func DeriveKey(in []byte) []byte {
	hash := sha256.Sum256(in)
	return hash[:]
}

func Padding(src []byte, blksize int) []byte {
	if blksize < 1 || blksize > 255 {
		panic(errInvalidPadding)
	}
	padding := blksize - len(src)%blksize
	padText := bytes.Repeat([]byte{paddingByte(padding)}, padding)
	out := make([]byte, len(padText)+len(src))
	copy(out, src)
	copy(out[len(src):], padText)
	return out
}

func paddingByte(padding int) byte {
	if padding < 0 || padding > 255 {
		panic(errInvalidPadding)
	}
	return byte(padding)
}

func UnPadding(src []byte, blksize int) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errInvalidPadding
	}
	padding := int(src[length-1])
	if padding > length || padding > blksize {
		return nil, errInvalidPadding
	}
	return src[:length-padding], nil
}

func Nonce(sz int) ([]byte, error) {
	nonce := make([]byte, sz)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("create nonce: %w", err)
	}
	return nonce, nil
}
