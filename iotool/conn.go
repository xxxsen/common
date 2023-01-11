package iotool

import (
	"io"
	"net"
)

type wrapConn struct {
	net.Conn
	r io.Reader
	w io.Writer
	c io.Closer
}

func WrapConn(basic net.Conn, r io.Reader, w io.Writer, c io.Closer) net.Conn {
	if r == nil {
		r = basic
	}
	if w == nil {
		w = basic
	}
	if c == nil {
		c = basic
	}
	return &wrapConn{
		Conn: basic,
		r:    r,
		w:    w,
		c:    c,
	}
}

func (w *wrapConn) Read(b []byte) (int, error) {
	return w.r.Read(b)
}

func (w *wrapConn) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *wrapConn) Close() error {
	return w.c.Close()
}
