package iotool

import (
	"net"
)

type StatisticConn struct {
	net.Conn
	rbytes uint64
	wbytes uint64
}

func NewStatisticConn(c net.Conn) *StatisticConn {
	return &StatisticConn{
		Conn:   c,
		rbytes: 0,
		wbytes: 0,
	}
}

func (s *StatisticConn) Read(b []byte) (int, error) {
	cnt, err := s.Conn.Read(b)
	if cnt > 0 {
		s.rbytes += uint64(cnt)
	}
	return cnt, err
}

func (s *StatisticConn) Write(b []byte) (int, error) {
	cnt, err := s.Conn.Write(b)
	if cnt > 0 {
		s.wbytes += uint64(cnt)
	}
	return cnt, err
}

func (s *StatisticConn) GetRBytes() uint64 {
	return s.rbytes
}

func (s *StatisticConn) GetWBytes() uint64 {
	return s.wbytes
}
