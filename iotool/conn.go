package iotool

import (
	"io"
	"net"
)

type IOConn struct {
	r     io.Reader
	w     io.Writer
	c     io.Closer
	lAddr net.Addr
	rAddr net.Addr
	net.Conn
}

type netAddr struct {
	network string
	addr    string
}

func (na netAddr) Network() string {
	return na.network
}

func (na netAddr) String() string {
	return na.addr
}

func newNetAddrWithString(network string, addr string) net.Addr {
	return netAddr{network: network, addr: addr}
}

func WrapConn(basic net.Conn, r io.Reader, w io.Writer, c io.Closer) net.Conn {
	ioc := NewIOConn(basic)
	if r != nil {
		ioc.ReplaceReader(r)
	}
	if w != nil {
		ioc.ReplaceWriter(w)
	}
	if c != nil {
		ioc.ReplaceCloser(c)
	}
	return ioc
}

func NewIOConn(base net.Conn) *IOConn {
	return &IOConn{
		r:     base,
		w:     base,
		c:     base,
		lAddr: base.LocalAddr(),
		rAddr: base.RemoteAddr(),
	}
}

func (c *IOConn) ReplaceReader(r io.Reader) *IOConn {
	c.r = r
	return c
}

func (c *IOConn) ReplaceWriter(w io.Writer) *IOConn {
	c.w = w
	return c
}

func (c *IOConn) ReplaceCloser(cr io.Closer) *IOConn {
	c.c = cr
	return c
}

func (c *IOConn) ReplaceLocalAddr(addr net.Addr) *IOConn {
	c.lAddr = addr
	return c
}

func (c *IOConn) ReplaceRemoteAddr(addr net.Addr) *IOConn {
	c.rAddr = addr
	return c
}

func (c *IOConn) ReplaceLocalAddrWithString(network string, addr string) *IOConn {
	return c.ReplaceLocalAddr(newNetAddrWithString(network, addr))
}

func (c *IOConn) ReplaceRemoteAddrWithString(network string, addr string) *IOConn {
	return c.ReplaceRemoteAddr(newNetAddrWithString(network, addr))
}

func (c *IOConn) Read(b []byte) (n int, err error) {
	return c.r.Read(b)
}

func (c *IOConn) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
}

func (c *IOConn) Close() error {
	return c.c.Close()
}

func (c *IOConn) LocalAddr() net.Addr {
	return c.lAddr
}

func (c *IOConn) RemoteAddr() net.Addr {
	return c.rAddr
}
