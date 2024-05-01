package httprewrite

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHTTPInfo(t *testing.T) {
	body := []byte("12345")
	r, err := http.NewRequest(http.MethodPost, "https://google.com?a=1&b=2", bytes.NewReader(body))
	assert.NoError(t, err)
	buf := bytes.NewBuffer(nil)
	err = r.WriteProxy(buf)
	assert.NoError(t, err)
	bc := BasicHTTPRequestContext{}
	bio := bufio.NewReader(buf)
	err = bc.Parse(bio)
	assert.NoError(t, err)
	assert.Equal(t, 1, bc.HTTPVersionMajor)
	assert.Equal(t, 1, bc.HTTPVersionMinor)
	assert.Equal(t, http.MethodPost, bc.Method)
	assert.Equal(t, "google.com", bc.URL.Host)
	assert.Equal(t, "1", bc.URL.Query().Get("a"))
	assert.Equal(t, "2", bc.URL.Query().Get("b"))
	data, _, err := bio.ReadLine()
	assert.NoError(t, err)
	assert.Equal(t, body, data)
	_, _, err = bio.ReadLine()
	assert.Equal(t, io.EOF, err)
}

func TestWriteNoProxy(t *testing.T) {
	body := []byte("12345")
	r, err := http.NewRequest(http.MethodPost, "https://google.com?a=1&b=2", bytes.NewReader(body))
	assert.NoError(t, err)
	buf := bytes.NewBuffer(nil)
	err = r.Write(buf)
	assert.NoError(t, err)
	bc := BasicHTTPRequestContext{}
	bio := bufio.NewReader(buf)
	err = bc.Parse(bio)
	assert.NoError(t, err)
	reader, err := bc.ToReader(false)
	assert.NoError(t, err)
	reader = io.MultiReader(reader, bio)
	bio = bufio.NewReader(reader)
	r, err = http.ReadRequest(bio)
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, data, body)
	assert.Equal(t, "google.com", r.Host)
	assert.Equal(t, "1", r.URL.Query().Get("a"))
	assert.Equal(t, "2", r.URL.Query().Get("b"))
}

func TestHeaderWrite(t *testing.T) {
	h := http.Header{}
	h.Set("a", "b")
	buf := bytes.NewBuffer(nil)
	err := h.Write(buf)
	assert.NoError(t, err)
	t.Logf(buf.String())
}
