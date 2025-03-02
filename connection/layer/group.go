package layer

import (
	"context"
	"fmt"
	"net"
)

type Group struct {
	ds []ILayer
}

func NewGroup(ds ...ILayer) ILayer {
	return &Group{
		ds: ds,
	}
}

func (g *Group) Name() string {
	return fmt.Sprintf("group:%+v", g.names(g.ds))
}

func (g *Group) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	var err error
	for _, d := range g.ds {
		conn, err = d.MakeLayerContext(ctx, conn)
		if err != nil {
			return nil, fmt.Errorf("make layer fail, layer name:%s, err:%w", d.Name(), err)
		}
	}
	return conn, nil
}

func (g *Group) names(ls []ILayer) []string {
	rs := make([]string, 0, len(ls))
	for _, l := range ls {
		rs = append(rs, l.Name())
	}
	return rs
}
