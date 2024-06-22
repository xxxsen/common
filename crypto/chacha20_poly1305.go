package crypto

import (
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

type chacha20Poly1305 struct {
}

func (c *chacha20Poly1305) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	nonce, err := Nonce(chacha20poly1305.NonceSize)
	if err != nil {
		return nil, err
	}
	ciphertext := block.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (c *chacha20Poly1305) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < chacha20poly1305.NonceSize {
		return nil, fmt.Errorf("invalid chipher text size:%d", len(ciphertext))
	}
	nonce, ciphertext := ciphertext[:chacha20poly1305.NonceSize], ciphertext[chacha20poly1305.NonceSize:]
	plaintext, err := block.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (c *chacha20Poly1305) Name() string {
	return CodecAesChacha20Poly1305
}

func NewChacha20Poly1305() ICodec {
	return &chacha20Poly1305{}
}

func init() {
	register(NewChacha20Poly1305())
}
