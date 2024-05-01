package padding

import (
	"context"
	"net"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/iotool"
	"github.com/xxxsen/common/utils"

	"github.com/xxxsen/common/errs"
)

func init() {
	layer.Regist(paddingLayerName, createPaddingLayer)
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
		return nil, errs.New(errs.ErrParam, "padding max < padding min")
	}
	return &paddingLayer{c: c}, nil
}

type paddingLayer struct {
	c *config
}

func (p *paddingLayer) Name() string {
	return paddingLayerName
}

func (d *paddingLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	c := iotool.WrapPadding(conn, d.c.PaddingMin, d.c.PaddingMax, d.c.PaddingIfLessThan, d.c.MaxBusiDataPerPacket)
	return iotool.WrapConn(conn, c, c, c), nil
}
