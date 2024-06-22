package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesCBC(t *testing.T) {
	data := []byte("hello world")
	key := DeriveKey([]byte("abc"))
	cbc := NewAesCBC()
	enc, err := cbc.Encrypt(data, key)
	assert.NoError(t, err)
	dec, err := cbc.Decrypt(enc, key)
	assert.NoError(t, err)
	assert.Equal(t, data, dec)
}
