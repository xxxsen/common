package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type aesStreamMode uint8

const (
	aesStreamCTR aesStreamMode = iota
	aesStreamCFB
	aesStreamOFB
)

func encryptAESStream(
	plaintext, key []byte,
	mode aesStreamMode,
) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	initializationVector := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, initializationVector); err != nil {
		return nil, fmt.Errorf("create initialization vector: %w", err)
	}
	stream := newAESStream(block, initializationVector, mode, false)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func decryptAESStream(
	ciphertext, key []byte,
	mode aesStreamMode,
) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errCiphertextTooShort
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}
	initializationVector := ciphertext[:aes.BlockSize]
	payload := ciphertext[aes.BlockSize:]
	stream := newAESStream(block, initializationVector, mode, true)
	stream.XORKeyStream(payload, payload)
	return payload, nil
}

func newAESStream(
	block cipher.Block,
	initializationVector []byte,
	mode aesStreamMode,
	decrypt bool,
) cipher.Stream {
	switch mode {
	case aesStreamCTR:
		return cipher.NewCTR(block, initializationVector)
	case aesStreamCFB:
		return newCFBStream(block, initializationVector, decrypt)
	case aesStreamOFB:
		return newOFBStream(block, initializationVector)
	}
	panic("unknown AES stream mode")
}

type cfbStream struct {
	block    cipher.Block
	state    []byte
	keyBlock []byte
	feedback []byte
	used     int
	decrypt  bool
}

func newCFBStream(block cipher.Block, initializationVector []byte, decrypt bool) cipher.Stream {
	size := block.BlockSize()
	stream := &cfbStream{
		block:    block,
		state:    append([]byte(nil), initializationVector...),
		keyBlock: make([]byte, size),
		feedback: make([]byte, size),
		decrypt:  decrypt,
	}
	block.Encrypt(stream.keyBlock, stream.state)
	return stream
}

func (s *cfbStream) XORKeyStream(destination, source []byte) {
	if len(destination) < len(source) {
		panic("crypto/cipher: output smaller than input")
	}
	for index, value := range source {
		if s.used == len(s.keyBlock) {
			copy(s.state, s.feedback)
			s.block.Encrypt(s.keyBlock, s.state)
			s.used = 0
		}
		output := value ^ s.keyBlock[s.used]
		destination[index] = output
		if s.decrypt {
			s.feedback[s.used] = value
		} else {
			s.feedback[s.used] = output
		}
		s.used++
	}
}

type ofbStream struct {
	block    cipher.Block
	state    []byte
	keyBlock []byte
	used     int
}

func newOFBStream(block cipher.Block, initializationVector []byte) cipher.Stream {
	size := block.BlockSize()
	return &ofbStream{
		block:    block,
		state:    append([]byte(nil), initializationVector...),
		keyBlock: make([]byte, size),
		used:     size,
	}
}

func (s *ofbStream) XORKeyStream(destination, source []byte) {
	if len(destination) < len(source) {
		panic("crypto/cipher: output smaller than input")
	}
	for index, value := range source {
		if s.used == len(s.keyBlock) {
			s.block.Encrypt(s.keyBlock, s.state)
			copy(s.state, s.keyBlock)
			s.used = 0
		}
		destination[index] = value ^ s.keyBlock[s.used]
		s.used++
	}
}
