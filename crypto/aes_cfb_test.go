package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesCrypto(t *testing.T) {
	cc := &aesCFB{}
	data := []byte("hello world")
	key := DeriveKey([]byte("abc123"))
	enc, err := cc.Encrypt(data, key)
	assert.NoError(t, err)
	dec, err := cc.Decrypt(enc, key)
	assert.NoError(t, err)
	assert.Equal(t, data, dec)
}
