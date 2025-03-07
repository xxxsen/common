package fragment

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/utils"
)

const (
	fragmentLayerName = "fragment"
)

type fragmentLayer struct {
	c *config
}

func (f *fragmentLayer) Name() string {
	return fragmentLayerName
}

func (f *fragmentLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	return newFragmentConn(conn, f.c), nil
}

type fragmentConn struct {
	isNeedFrag bool
	c          *config
	packetId   uint32
	net.Conn
}

func newFragmentConn(conn net.Conn, c *config) net.Conn {
	return &fragmentConn{Conn: conn, c: c, isNeedFrag: true, packetId: 0}
}

func (c *fragmentConn) checkNumberMatch(in []uint32, num uint32) (bool, bool) {
	if len(in) == 0 {
		return false, false
	}
	var match bool
	var checkNextTime bool
	if len(in) == 1 {
		match = in[0] == num
		checkNextTime = num < in[0]
		return match, checkNextTime
	}
	match = in[0] <= num && in[1] >= num
	checkNextTime = num < in[1]
	return match, checkNextTime
}

func (c *fragmentConn) calcRandRange(rng []uint32) uint32 {
	if len(rng) == 0 {
		return 0
	}
	left := rng[0]
	right := left
	if len(rng) > 1 {
		right = rng[1]
	}
	if right == left {
		return right
	}
	return left + rand.Uint32()%(right-left)
}

func (c *fragmentConn) calcNextInterval(interval []uint32) uint32 {
	return c.calcRandRange(interval)
}

func (c *fragmentConn) calcNextBytesToSend(b []uint32) uint32 {
	if v := c.calcRandRange(b); v != 0 {
		return v
	}
	return math.MaxUint16

}

func (c *fragmentConn) Write(b []byte) (int, error) {
	if !c.isNeedFrag {
		return c.Conn.Write(b)
	}
	match, checkNextTime := c.checkNumberMatch(c.c.PacketNumberRange, c.packetId)
	c.packetId++
	if !checkNextTime {
		c.isNeedFrag = false
	}
	if !match {
		return c.Conn.Write(b)
	}
	sz := len(b)
	for len(b) > 0 {
		nextToSend := c.calcNextBytesToSend(c.c.PacketLengthRange)
		if nextToSend > uint32(len(b)) {
			nextToSend = uint32(len(b))
		}
		buf := b[:nextToSend]
		b = b[nextToSend:]
		if _, err := c.Conn.Write(buf); err != nil {
			return 0, fmt.Errorf("send part pkt failed, err:%w", err)
		}
		nextInterval := c.calcNextInterval(c.c.IntervalRange)
		if len(b) != 0 && nextInterval != 0 {
			time.Sleep(time.Duration(nextInterval) * time.Millisecond)
		}
	}
	return sz, nil
}

func createFragmentLayer(param interface{}) (layer.ILayer, error) {
	dst := &config{}
	if err := utils.ConvStructJson(param, dst); err != nil {
		return nil, err
	}
	if len(dst.IntervalRange) > 2 || len(dst.PacketLengthRange) > 2 || len(dst.PacketNumberRange) > 2 {
		return nil, fmt.Errorf("invalid fragment params")
	}
	if len(dst.IntervalRange) == 2 && dst.IntervalRange[0] > dst.IntervalRange[1] {
		return nil, fmt.Errorf("invalid interval range")
	}
	if len(dst.PacketLengthRange) == 2 && dst.PacketLengthRange[0] > dst.PacketLengthRange[1] {
		return nil, fmt.Errorf("invalid packet length range")
	}
	if len(dst.PacketNumberRange) == 2 && dst.PacketNumberRange[0] > dst.PacketNumberRange[1] {
		return nil, fmt.Errorf("invalid packet number range")
	}
	return &fragmentLayer{c: dst}, nil
}

func init() {
	layer.Register(fragmentLayerName, createFragmentLayer)
}
