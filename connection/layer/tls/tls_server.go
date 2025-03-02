package tls

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/utils"
)

func init() {
	layer.Register(tlsServerLayerName, createTlsServerLayer)
}

func createTlsServerLayer(params interface{}) (layer.ILayer, error) {
	c := &svrConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	cer, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("invalid cert info, err:%w", err)
	}
	return &tlsServerLayer{
		c: c,
		tlsc: &tls.Config{
			Certificates: []tls.Certificate{cer},
			MinVersion:   GetVersionIdByNameOrDefault(c.MinTLSVersion),
			MaxVersion:   GetVersionIdByNameOrDefault(c.MaxTLSVersion),
		},
	}, nil
}

type tlsServerLayer struct {
	c    *svrConfig
	tlsc *tls.Config
}

func (d *tlsServerLayer) Name() string {
	return tlsServerLayerName
}

func (d *tlsServerLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	tlsconn := tls.Server(conn, d.tlsc)
	if err := tlsconn.HandshakeContext(ctx); err != nil {
		return nil, fmt.Errorf("handshake tls failed, err:%w", err)
	}
	return tlsconn, nil
}
