package tls

import (
	"context"
	ctls "crypto/tls"
	"fmt"
	"math/rand"
	"net"

	"github.com/xxxsen/common/connection/layer"

	utls "github.com/refraction-networking/utls"

	"github.com/xxxsen/common/utils"
)

func init() {
	layer.Register(tlsDialLayerName, createTLSDialLayer)
}

type itlsconn interface {
	net.Conn
	VerifyHostname(host string) error
	HandshakeContext(ctx context.Context) error
}

func createTLSDialLayer(params interface{}) (layer.ILayer, error) {
	c := &cliConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	if len(c.FingerPrint) > 0 {
		if _, ok := mappingHelloID[c.FingerPrint]; !ok {
			return nil, fmt.Errorf("invalid finger print str:%s", c.FingerPrint)
		}
	}
	if len(c.SNI) != 0 {
		c.SNIs = append(c.SNIs, c.SNI)
	}
	c.SNIs = dedup(c.SNIs)
	return &tlsDialLayer{c: c}, nil
}

func dedup(lst []string) []string {
	rs := make([]string, 0, len(lst))
	exist := make(map[string]struct{}, len(lst))
	for _, item := range lst {
		if _, ok := exist[item]; ok {
			continue
		}
		rs = append(rs, item)
		exist[item] = struct{}{}
	}
	return rs
}

type tlsDialLayer struct {
	c *cliConfig
}

func (d *tlsDialLayer) Name() string {
	return tlsDialLayerName
}

func (d *tlsDialLayer) selectSNI() string {
	if len(d.c.SNIs) == 0 {
		return ""
	}
	return d.c.SNIs[rand.Int()%len(d.c.SNIs)]
}

func (d *tlsDialLayer) createTlsConn(conn net.Conn, fingerprint string) itlsconn {
	if len(d.c.FingerPrint) == 0 {
		tlsconn := ctls.Client(conn, &ctls.Config{
			ServerName:         d.selectSNI(),
			InsecureSkipVerify: d.c.SkipInsecureVerify,
			MinVersion:         GetVersionIdByNameOrDefault(d.c.MinTLSVersion),
			MaxVersion:         GetVersionIdByNameOrDefault(d.c.MaxTLSVersion),
		})
		return tlsconn
	}
	return utls.UClient(conn, &utls.Config{
		ServerName:         d.c.SNI,
		InsecureSkipVerify: d.c.SkipInsecureVerify,
	}, mappingHelloID[fingerprint])
}

func (d *tlsDialLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	tlsconn := d.createTlsConn(conn, d.c.FingerPrint)
	if err := tlsconn.HandshakeContext(ctx); err != nil {
		return nil, fmt.Errorf("handshake tls failed, err:%w", err)
	}
	if !d.c.SkipInsecureVerify {
		if err := tlsconn.VerifyHostname(d.c.SNI); err != nil {
			return nil, fmt.Errorf("verify tls host name failed, err:%w", err)
		}
	}
	return tlsconn, nil
}
