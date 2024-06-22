package crypto

import "crypto/sha256"

func DeriveKey(in []byte) []byte {
	hash := sha256.Sum256(in)
	return hash[:]
}
