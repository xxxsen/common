package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesGCM(t *testing.T) {
	data := []byte("hello world")
	key := DeriveKey([]byte("123"))
	gcm := NewAesGCM()
	enc, err := gcm.Encrypt(data, key)
	assert.NoError(t, err)
	dec, err := gcm.Decrypt(enc, key)
	assert.NoError(t, err)
	assert.Equal(t, data, dec)
}
