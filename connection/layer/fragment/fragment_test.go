package fragment

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeConn struct {
	net.Conn
	buf       bytes.Buffer
	sliceList [][]byte
}

func (c *fakeConn) Read(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *fakeConn) Write(b []byte) (n int, err error) {
	c.sliceList = append(c.sliceList, b)
	return c.buf.Write(b)
}

func (c *fakeConn) Close() error { return nil }

func (c *fakeConn) LocalAddr() net.Addr { return nil }

func (c *fakeConn) RemoteAddr() net.Addr { return nil }

func (c *fakeConn) SetDeadline(t time.Time) error { return nil }

func (c *fakeConn) SetReadDeadline(t time.Time) error { return nil }

func (c *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

func TestFragment(t *testing.T) {
	wordlist := []string{
		"hello world, adadaaf,af afd asfa faf ,adad afas asd",
		"hi, this is a test, fsdfsafa f,a faf ,a g 4 ,sef f,esf ",
		"r u ok? adsasj fja fjaf ajfa hfaf iaua;,fa asd",
	}
	const maxPacketSize = 10
	c := &config{
		IntervalRange:     []uint32{100, 200},
		PacketLengthRange: []uint32{1, maxPacketSize},
		PacketNumberRange: []uint32{0, 100},
	}
	ly, err := createFragmentLayer(c)
	assert.NoError(t, err)
	fconn := &fakeConn{}
	conn, err := ly.MakeLayerContext(context.Background(), fconn)
	assert.NoError(t, err)
	for i, data := range wordlist {
		log.Printf("write index:%d", i)
		_, err = conn.Write([]byte(data))
		assert.NoError(t, err)
	}
	raw, err := io.ReadAll(&fconn.buf)
	assert.NoError(t, err)
	wrote := strings.Join(wordlist, "")
	assert.Equal(t, string(raw), wrote)
	for _, b := range fconn.sliceList {
		assert.LessOrEqual(t, len(b), maxPacketSize)
	}
}
