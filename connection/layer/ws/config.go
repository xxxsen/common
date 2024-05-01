package ws

type cliConfig struct {
	Schema              string `json:"schema"`
	Host                string `json:"host"`
	Path                string `json:"path"`
	HandshakePaddingMin int    `json:"handshake_padding_min"`
	HandshakePaddingMax int    `json:"handshake_padding_max"`
}

type svrConfig struct {
	Path                string `json:"path"`
	Host                string `json:"host"`
	HandshakePaddingMin int    `json:"handshake_padding_min"`
	HandshakePaddingMax int    `json:"handshake_padding_max"`
}
