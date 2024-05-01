package httprewrite

import (
	"bufio"
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/xxxsen/common/iotool"

	"github.com/stretchr/testify/assert"
)

type fakeconn struct {
	net.TCPConn
}

func TestHTTPRewrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	r, err := http.NewRequest(http.MethodPut, "https://test.com?q1=q2", bytes.NewReader([]byte("12345")))
	assert.NoError(t, err)
	r.Header.Set("h1", "v1")
	err = r.Write(buf)
	assert.NoError(t, err)
	conn := iotool.WrapConn(&fakeconn{}, buf, nil, http.NoBody)

	dialer, err := createHTTPRequestRewriteLayer(&config{
		ForceUseProxy: true,
		RewritePath:   "/abc",
		RewriteQuery:  map[string]string{"q1": "q1-1"},
		RewriteHeader: map[string]string{"h1": "h1-1", "host": "new-test.com"},
	})
	assert.NoError(t, err)
	conn, err = dialer.MakeLayerContext(context.Background(), conn)
	assert.NoError(t, err)
	bio := bufio.NewReader(conn)
	r, err = http.ReadRequest(bio)
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPut, r.Method)
	assert.Equal(t, data, []byte("12345"))
	assert.Equal(t, "new-test.com", r.Host)
	assert.Equal(t, "/abc", r.URL.Path)
	assert.Equal(t, "q1-1", r.URL.Query().Get("q1"))
	assert.Equal(t, "h1-1", r.Header.Get("h1"))
}
