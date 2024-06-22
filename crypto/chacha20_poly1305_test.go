package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChacha20Pol1305(t *testing.T) {
	data := []byte("hello world")
	key := DeriveKey([]byte("123"))
	cc := NewChacha20Poly1305()
	enc, err := cc.Encrypt(data, key)
	assert.NoError(t, err)
	dec, err := cc.Decrypt(enc, key)
	assert.NoError(t, err)
	assert.Equal(t, data, dec)
}
