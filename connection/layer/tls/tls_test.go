package tls

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	"github.com/xxxsen/common/iotool"
)

func handleServer(c net.Conn) {
	ly, err := createTlsServerLayer(&svrConfig{
		CertFile: "/tmp/server.crt",
		KeyFile:  "/tmp/server.key",
	})
	if err != nil {
		panic(err)
	}
	c, err = ly.MakeLayerContext(context.Background(), c)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 64)
	_, err = io.ReadAtLeast(c, buf, len(buf))
	if err != nil {
		panic(err)
	}
	log.Printf("recv msg:%s\n", string(buf))
}

func handleClient(c net.Conn) {
	ly, err := createTLSDialLayer(&cliConfig{
		SNI:                "abc.com",
		SkipInsecureVerify: false,
	})
	if err != nil {
		panic(err)
	}
	c, err = ly.MakeLayerContext(context.Background(), c)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 64)
	copy(buf, []byte("hello server"))
	_, err = c.Write(buf)
	if err != nil {
		panic(err)
	}
}

func TestTLS(t *testing.T) {
	clientReader, serverWriter := io.Pipe()
	serverReader, clientWriter := io.Pipe()
	clientRw := iotool.WrapReadWriter(clientReader, clientWriter)
	serverRw := iotool.WrapReadWriter(serverReader, serverWriter)
	clientConn := iotool.WrapConn(nil, clientRw, clientRw, nil)
	serverConn := iotool.WrapConn(nil, serverRw, serverRw, nil)
	go func() {
		handleServer(serverConn)
	}()
	handleClient(clientConn)

}
