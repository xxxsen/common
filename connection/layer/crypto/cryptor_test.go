package crypto

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/xxxsen/common/crypto"
	"github.com/xxxsen/common/iotool"
)

type simplerwc struct {
	bytes.Buffer
}

func (s *simplerwc) Close() error {
	return nil
}

func TestCryptor(t *testing.T) {
	data := make([]byte, 513*1024)
	io.ReadFull(rand.Reader, data)
	ly, err := createCryptor(&config{
		Key:   "hello world",
		Codec: "aes-gcm",
	})
	assert.NoError(t, err)
	buf := &simplerwc{}
	conn := net.Conn(iotool.WrapReadWriteCloserToIOConn(buf))
	conn, err = ly.MakeLayerContext(context.Background(), conn)
	assert.NoError(t, err)
	_, err = conn.Write(data)
	assert.NoError(t, err)
	assert.NotEqual(t, data, buf.Buffer.Bytes())
	readbuf := make([]byte, len(data))
	_, err = io.ReadAtLeast(conn, readbuf, len(readbuf))
	assert.NoError(t, err)
	assert.Equal(t, data, readbuf)
}
