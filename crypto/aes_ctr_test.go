package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesCtr(t *testing.T) {
	data := []byte("hello world")
	key := DeriveKey([]byte("abc"))
	ctr := NewAesCTR()
	enc, err := ctr.Encrypt(data, key)
	assert.NoError(t, err)
	dec, err := ctr.Decrypt(enc, key)
	assert.NoError(t, err)
	assert.Equal(t, data, dec)
}
