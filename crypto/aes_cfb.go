package crypto

type aesCFB struct {
}

func (c *aesCFB) Encrypt(buf []byte, key []byte) ([]byte, error) {
	return encryptAESStream(buf, key, aesStreamCFB)
}

func (c *aesCFB) Decrypt(buf []byte, key []byte) ([]byte, error) {
	return decryptAESStream(buf, key, aesStreamCFB)
}

func (c *aesCFB) Name() string {
	return CodecAesCFB
}

func NewAesCFB() ICodec {
	return &aesCFB{}
}
