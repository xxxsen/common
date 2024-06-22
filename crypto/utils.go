package crypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

func DeriveKey(in []byte) []byte {
	hash := sha256.Sum256(in)
	return hash[:]
}

func Padding(src []byte, blksize int) []byte {
	padding := blksize - len(src)%blksize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	out := make([]byte, len(padText)+len(src))
	copy(out, src)
	copy(out[len(src):], padText)
	return out
}

func UnPadding(src []byte, blksize int) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("unpad: input data is empty")
	}
	padding := int(src[length-1])
	if padding > length || padding > blksize {
		return nil, fmt.Errorf("unpad: invalid padding size")
	}
	return src[:length-padding], nil
}
