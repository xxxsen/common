package utils

import (
	"crypto/rand"
	"math/big"
)

func RandBytes(l, r int) []byte {
	if l > r {
		l = r
	}
	length := secureRandomInt(l, r)
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return buf
}

func secureRandomInt(minimum, maximum int) int {
	if minimum == maximum {
		return minimum
	}
	span := big.NewInt(int64(maximum - minimum + 1))
	offset, err := rand.Int(rand.Reader, span)
	if err != nil {
		panic(err)
	}
	return minimum + int(offset.Int64())
}
