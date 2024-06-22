package crypto

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWrite(t *testing.T) {
	ccs := []ICodec{
		NewAesCBC(),
		NewAesCFB(),
		NewAesCTR(),
		NewAesGCM(),
		NewAesOFB(),
		NewChacha20Poly1305(),
	}
	szs := []int{
		31*1024 + 123,
		34*1024 + 235,
		63 * 1024,
		63*1024 + 33,
		64*1024 + 55,
		136*1024 + 221,
	}
	for _, sz := range szs {
		for _, cc := range ccs {
			buf := make([]byte, sz)
			_, err := rand.Read(buf)
			savedBuf := make([]byte, len(buf))
			copy(savedBuf, buf)
			assert.NoError(t, err)
			key := DeriveKey([]byte("hello world"))
			conn := &bytes.Buffer{}
			writer := NewWriter(conn, cc, key)
			_, err = writer.Write(buf)
			assert.NoError(t, err)
			reader := NewReader(conn, cc, key)
			data, err := io.ReadAll(reader)
			assert.NoError(t, err)
			assert.Equal(t, savedBuf, data)
		}
	}
}
