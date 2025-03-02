package iotool

import (
	"fmt"
	"io"
	"net"
	"time"
)

type readOnlyConn struct {
	reader io.Reader
}

func NewReadOnlyConn(r io.Reader) net.Conn {
	return &readOnlyConn{
		reader: r,
	}
}

func (conn readOnlyConn) Read(p []byte) (int, error) { return conn.reader.Read(p) }
func (conn readOnlyConn) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("not allow to write")
}
func (conn readOnlyConn) Close() error                       { return nil }
func (conn readOnlyConn) LocalAddr() net.Addr                { return nil }
func (conn readOnlyConn) RemoteAddr() net.Addr               { return nil }
func (conn readOnlyConn) SetDeadline(t time.Time) error      { return nil }
func (conn readOnlyConn) SetReadDeadline(t time.Time) error  { return nil }
func (conn readOnlyConn) SetWriteDeadline(t time.Time) error { return nil }
