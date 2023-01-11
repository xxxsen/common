package iotool

import (
	"bufio"
	"net"
)

type bufRConn struct {
	r *bufio.Reader
	net.Conn
}

func WrapBufRConn(c net.Conn) net.Conn {
	return &bufRConn{
		r:    bufio.NewReader(c),
		Conn: c,
	}
}

func (c *bufRConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}
