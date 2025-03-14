package message

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/utils"
)

func init() {
	layer.Register(messageLayerName, createMessageLayer)
}

func createMessageLayer(params interface{}) (layer.ILayer, error) {
	c := &config{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	for idx, act := range c.Actions {
		if act.Length == 0 {
			return nil, fmt.Errorf("invalid length, idx:%d", idx)
		}
		if strings.EqualFold(act.Type, actionTypeSend) && strings.EqualFold(act.Type, actionTypeRecv) {
			return nil, fmt.Errorf("action type invalid")
		}
	}
	return &messageLayer{c: c}, nil
}

type messageLayer struct {
	c *config
}

func (d *messageLayer) Name() string {
	return messageLayerName
}

func (d *messageLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	for idx, act := range d.c.Actions {
		if strings.EqualFold(act.Type, actionTypeSend) {
			if _, err := conn.Write(utils.RandBytes(int(act.Length), int(act.Length))); err != nil {
				return nil, fmt.Errorf("send msg fail at idx:%d, err:%w", idx, err)
			}
		}
		if strings.EqualFold(act.Type, actionTypeRecv) {
			if _, err := io.CopyN(io.Discard, conn, int64(act.Length)); err != nil {
				return nil, fmt.Errorf("recv msg fail at idx:%d, err:%w", idx, err)
			}
		}
	}
	return conn, nil
}
