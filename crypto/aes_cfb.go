package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type aesCFB struct {
}

func (c *aesCFB) Encrypt(buf []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(buf))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], buf)

	return ciphertext, nil
}

func (c *aesCFB) Decrypt(buf []byte, key []byte) ([]byte, error) {
	if len(buf) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(buf, buf)

	return buf, nil
}

func NewAesCFB() ICodec {
	return &aesCFB{}
}
