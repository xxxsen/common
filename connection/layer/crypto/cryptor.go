package crypto

import (
	"context"
	"fmt"
	"net"

	"github.com/xxxsen/common/connection/layer"
	"github.com/xxxsen/common/crypto"
	"github.com/xxxsen/common/iotool"
	"github.com/xxxsen/common/utils"
)

type cryptor struct {
	c  *config
	cc crypto.ICodec
}

func (c *cryptor) Name() string {
	return encryptorLayerName
}

func (c *cryptor) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	key := crypto.DeriveKey([]byte(c.c.Key))
	r := crypto.NewReader(conn, c.cc, key)
	w := crypto.NewWriter(conn, c.cc, key)
	return iotool.WrapConn(conn, r, w, nil), nil
}

func createCryptor(params interface{}) (layer.ILayer, error) {
	c := &config{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	cc, ok := crypto.FindCodec(c.Codec)
	if !ok {
		return nil, fmt.Errorf("codec:%s not found", c.Codec)
	}

	return &cryptor{cc: cc, c: c}, nil
}

func init() {
	layer.Register("cryptor", createCryptor)
}
