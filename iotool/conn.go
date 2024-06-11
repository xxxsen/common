package iotool

import (
	"fmt"
	"io"
	"net"
	"time"
)

type setDeadLineFunc func(t time.Time) error

func defaultSetDeadline(t time.Time) error {
	return fmt.Errorf("not impl")
}

var defaultNoAddr = newNetAddrWithString("unknown", "no addr")

type IOConn struct {
	r                io.Reader
	w                io.Writer
	c                io.Closer
	lAddr            net.Addr
	rAddr            net.Addr
	setReadDeadLine  setDeadLineFunc
	setWriteDeadLine setDeadLineFunc
	setDeadLine      setDeadLineFunc
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
	return WrapConnToIOConn(basic, r, w, c)
}

func WrapConnToIOConn(basic net.Conn, r io.Reader, w io.Writer, c io.Closer) *IOConn {
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

func WrapReadWriteCloserToIOConn(rwc io.ReadWriteCloser) *IOConn {
	return newIOConn(rwc, rwc, rwc, defaultNoAddr, defaultNoAddr, defaultSetDeadline, defaultSetDeadline, defaultSetDeadline)
}

func NewIOConn(base net.Conn) *IOConn {
	return newIOConn(base, base, base, base.LocalAddr(), base.RemoteAddr(), base.SetDeadline, base.SetReadDeadline, base.SetWriteDeadline)
}

func newIOConn(r io.Reader, w io.Writer, c io.Closer, laddr net.Addr, raddr net.Addr, sd, srd, swd setDeadLineFunc) *IOConn {
	return &IOConn{
		r:                r,
		w:                w,
		c:                c,
		lAddr:            laddr,
		rAddr:            raddr,
		setDeadLine:      sd,
		setReadDeadLine:  srd,
		setWriteDeadLine: swd,
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

func (c *IOConn) SetDeadline(t time.Time) error {
	return c.setDeadLine(t)
}

func (c *IOConn) SetReadDeadline(t time.Time) error {
	return c.setReadDeadLine(t)
}

func (c *IOConn) SetWriteDeadline(t time.Time) error {
	return c.setWriteDeadLine(t)
}
