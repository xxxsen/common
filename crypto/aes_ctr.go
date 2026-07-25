package crypto

type aesCTR struct {
}

func (c *aesCTR) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	return encryptAESStream(plaintext, key, aesStreamCTR)
}

func (c *aesCTR) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	return decryptAESStream(ciphertext, key, aesStreamCTR)
}

func (c *aesCTR) Name() string {
	return CodecAesCTR
}

func NewAesCTR() ICodec {
	return &aesCTR{}
}
