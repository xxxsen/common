package crypto

type ICodec interface {
	Name() string
	Encrypt(buf []byte, key []byte) ([]byte, error)
	Decrypt(buf []byte, key []byte) ([]byte, error)
}

var codecMapping = make(map[string]ICodec)

func register(cc ICodec) {
	codecMapping[cc.Name()] = cc
}

func FindCodec(name string) (ICodec, bool) {
	if c, ok := codecMapping[name]; ok {
		return c, true
	}
	return nil, false
}
