package httprewrite

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

type BasicHTTPRequestContext struct {
	Method           string
	URL              *url.URL
	HTTPVersionMajor int
	HTTPVersionMinor int
	Header           http.Header
}

func ParseBasicHTTPRequestContext(bio *bufio.Reader) (*BasicHTTPRequestContext, error) {
	c := &BasicHTTPRequestContext{}
	if err := c.Parse(bio); err != nil {
		return nil, fmt.Errorf("parse failed, err:%w", err)
	}
	return c, nil
}

func (c *BasicHTTPRequestContext) Parse(bio *bufio.Reader) error {
	reader := textproto.NewReader(bio)
	//GET /xxx?a=b&x=y HTTP/1.1\r\n
	//H1: V1\r\n
	//H2: V2\r\n
	//\r\n
	firstline, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("read first line failed, err:%w", err)
	}
	header, err := reader.ReadMIMEHeader()
	if err != nil {
		return fmt.Errorf("read header failed, err:%w", err)
	}
	c.Header = http.Header(header)
	method, rest, ok1 := strings.Cut(firstline, " ")
	requestURI, proto, ok2 := strings.Cut(rest, " ")
	if !ok1 || !ok2 {
		return fmt.Errorf("invalid http first line:<%s>", hex.EncodeToString([]byte(firstline)))
	}
	justAuthority := strings.EqualFold(method, http.MethodConnect) && !strings.HasPrefix(requestURI, "/")
	if justAuthority {
		requestURI = "http://" + requestURI
	}
	uri, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return fmt.Errorf("parse request uri failed, uri:%s, err:%w", requestURI, err)
	}
	mj, mn, ok := http.ParseHTTPVersion(proto)
	if !ok {
		return fmt.Errorf("parse http version failed, err:%w", err)
	}
	c.Method = method
	c.URL = uri
	c.HTTPVersionMajor = mj
	c.HTTPVersionMinor = mn
	return nil

}

func (c *BasicHTTPRequestContext) buildUri(usingProxy bool) (string, error) {
	ruri := c.URL.RequestURI()
	host := c.URL.Host
	if len(host) == 0 {
		host = c.Header.Get("host")
	}
	if usingProxy && c.URL.Scheme != "" && c.URL.Opaque == "" {
		ruri = c.URL.Scheme + "://" + host + ruri
	} else if c.Method == "CONNECT" && c.URL.Path == "" {
		ruri = host
		if c.URL.Opaque != "" {
			ruri = c.URL.Opaque
		}
	}
	return ruri, nil
}

func (c *BasicHTTPRequestContext) ToReader(useproxy bool) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	ruri, err := c.buildUri(useproxy)
	if err != nil {
		return nil, fmt.Errorf("build ruri failed, err:%w", err)
	}
	buf.WriteString(fmt.Sprintf("%s %s HTTP/%d.%d\r\n", c.Method, ruri, c.HTTPVersionMajor, c.HTTPVersionMinor))
	if err := c.Header.Write(buf); err != nil {
		return nil, fmt.Errorf("write header failed, err:%w", err)
	}
	_, _ = io.WriteString(buf, "\r\n")
	return buf, nil
}
