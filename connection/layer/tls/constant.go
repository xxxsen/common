package tls

import (
	utls "github.com/refraction-networking/utls"
)

const (
	tlsDialLayerName   = "tls_client"
	tlsServerLayerName = "tls_server"
)

var mappingHelloID = map[string]utls.ClientHelloID{
	"chrome":  utls.HelloChrome_Auto,
	"360":     utls.Hello360_Auto,
	"firefox": utls.HelloFirefox_Auto,
	"safari":  utls.HelloSafari_Auto,
	"ios":     utls.HelloIOS_Auto,
	"android": utls.HelloAndroid_11_OkHttp,
	"qq":      utls.HelloQQ_Auto,
	"edge":    utls.HelloEdge_Auto,
	"random":  utls.HelloRandomized,
	"golang":  utls.HelloGolang,
}
