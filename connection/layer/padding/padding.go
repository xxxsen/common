package padding

import (
	"context"
	"fmt"
	"net"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/iotool"
	"github.com/xxxsen/common/utils"
)

func init() {
	layer.Register(paddingLayerName, createPaddingLayer)
}

func createPaddingLayer(param interface{}) (layer.ILayer, error) {
	c := &config{}
	if err := utils.ConvStructJson(param, c); err != nil {
		return nil, err
	}
	if c.PaddingMin == 0 {
		c.PaddingMin = 200
	}
	if c.PaddingMax == 0 {
		c.PaddingMax = 1000
	}
	if c.PaddingIfLessThan == 0 {
		c.PaddingIfLessThan = 4096
	}
	if c.MaxBusiDataPerPacket == 0 {
		c.MaxBusiDataPerPacket = 16384
	}
	if c.PaddingMax < c.PaddingMin {
		return nil, fmt.Errorf("padding max < padding min")
	}
	return &paddingLayer{c: c}, nil
}

type paddingLayer struct {
	c *config
}

func (p *paddingLayer) Name() string {
	return paddingLayerName
}

func (p *paddingLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	c := iotool.WrapPadding(conn, p.c.PaddingMin, p.c.PaddingMax, p.c.PaddingIfLessThan, p.c.MaxBusiDataPerPacket)
	return iotool.WrapConn(conn, c, c, c), nil
}
