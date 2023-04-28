package utils

import (
	"crypto/rand"
	unsaferand "math/rand"
)

//import

func RandBytes(l, r int) []byte {
	if l > r {
		l = r
	}
	length := unsaferand.Intn(r-l+1) + l
	buf := make([]byte, length)
	_, _ = rand.Reader.Read(buf)
	return buf 
}
