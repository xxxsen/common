package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesOFB(t *testing.T) {
	data := []byte("hello world")
	key := DeriveKey([]byte("123"))
	ofb := NewAesOFB()
	enc, err := ofb.Encrypt(data, key)
	assert.NoError(t, err)
	dec, err := ofb.Decrypt(enc, key)
	assert.NoError(t, err)
	assert.Equal(t, data, dec)
}
