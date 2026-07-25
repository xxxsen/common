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
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create AES-GCM: %w", err)
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
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create AES-GCM: %w", err)
	}
	if len(buf) < aesGCM.NonceSize() {
		return nil, errCiphertextTooShort
	}
	nonce, buf := buf[:aesGCM.NonceSize()], buf[aesGCM.NonceSize():]
	plaintext, err := aesGCM.Open(nil, nonce, buf, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt AES-GCM: %w", err)
	}

	return plaintext, nil
}

func (c *aesGCM) Name() string {
	return CodecAesGCM
}

func NewAesGCM() ICodec {
	return &aesGCM{}
}
