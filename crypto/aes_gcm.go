package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type aesGCM struct {
}

func (c *aesGCM) Encrypt(buf []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce, err := Nonce(aesGCM.NonceSize())
	if err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, buf, nil)
	return ciphertext, nil
}

func (c *aesGCM) Decrypt(buf []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(buf) < aesGCM.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, buf := buf[:aesGCM.NonceSize()], buf[aesGCM.NonceSize():]
	plaintext, err := aesGCM.Open(nil, nonce, buf, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c *aesGCM) Name() string {
	return CodecAesGCM
}

func NewAesGCM() ICodec {
	return &aesGCM{}
}

func init() {
	register(NewAesGCM())
}
