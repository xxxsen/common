package crypto

type aesOFB struct {
}

func (c *aesOFB) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	return encryptAESStream(plaintext, key, aesStreamOFB)
}

func (c *aesOFB) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	return decryptAESStream(ciphertext, key, aesStreamOFB)
}

func (c *aesOFB) Name() string {
	return CodecAesOFB
}

func NewAesOFB() ICodec {
	return &aesOFB{}
}
