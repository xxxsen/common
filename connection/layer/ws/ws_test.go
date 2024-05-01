package ws

import (
	"context"
	"io"
	"log"
	"net"
	"testing"

	"github.com/xxxsen/common/iotool"
)

func handleServer(c net.Conn) {
	ly, err := createWsServerLayer(&svrConfig{
		Path:                "/abc",
		HandshakePaddingMin: 20,
		HandshakePaddingMax: 1000,
	})
	if err != nil {
		panic(err)
	}
	c, err = ly.MakeLayerContext(context.Background(), c)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 64)
	copy(buf, []byte("hello client"))
	_, err = c.Write(buf)
	if err != nil {
		panic(err)
	}
	_, err = io.ReadAtLeast(c, buf, len(buf))
	if err != nil {
		panic(err)
	}
	log.Printf("server read:%s\n", string(buf))
}

func handleClient(c net.Conn) {
	ly, err := createWsDialLayer(&cliConfig{
		Schema:              "http",
		Host:                "abc.com",
		Path:                "/abc",
		HandshakePaddingMin: 100,
		HandshakePaddingMax: 550,
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
	log.Printf("client read:%s\n", string(buf))
	copy(buf, []byte("hello server"))
	_, err = c.Write(buf)
	if err != nil {
		panic(err)
	}
}

func TestReadWrite(t *testing.T) {
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
