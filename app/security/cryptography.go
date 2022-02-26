package security

import (
	"TeachAssistApi/app"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
)

type CryptographyService struct {
	gcm cipher.AEAD
}

func NewCryptographyService() (cs *CryptographyService, err error) {
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, app.CreateError(app.CryptographyError)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, app.CreateError(app.CryptographyError)
	}
	return &CryptographyService{gcm: gcm}, nil
}

func (cs CryptographyService) Encrypt(plain []byte) (cipher []byte, err error) {
	nonce := make([]byte, cs.gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, app.CreateError(app.CryptographyError)
	}
	cipher = cs.gcm.Seal(nil, nonce, plain, nil)
	cipher = append(nonce, cipher...)
	return cipher, nil
}

func (cs CryptographyService) EncryptBase64(plain []byte) ([]byte, error) {
	c, err := cs.Encrypt(plain)
	if err != nil {
		return nil, app.CreateError(app.CryptographyError)
	}

	base64Cipher := make([]byte, base64.RawStdEncoding.EncodedLen(len(c)))
	base64.RawStdEncoding.Encode(base64Cipher, c)

	return base64Cipher, nil
}

func (cs CryptographyService) Decrypt(cipher []byte) ([]byte, error) {
	nonce := cipher[0:cs.gcm.NonceSize()]
	ciphertext := cipher[cs.gcm.NonceSize():]

	plain, err := cs.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, app.CreateError(app.CryptographyError)
	}

	return plain, nil
}

func (cs CryptographyService) DecryptBase64(base64Cipher []byte) ([]byte, error) {
	c := make([]byte, base64.StdEncoding.DecodedLen(len(base64Cipher)))
	_, err := base64.StdEncoding.Decode(c, base64Cipher)
	if err != nil {
		return nil, app.CreateError(app.CryptographyError)
	}

	return cs.Decrypt(c)
}
