package crypto

type ICodec interface {
	Encrypt(buf []byte, key []byte) ([]byte, error)
	Decrypt(buf []byte, key []byte) ([]byte, error)
}
