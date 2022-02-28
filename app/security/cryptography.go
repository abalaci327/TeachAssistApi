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

func (cs CryptographyService) EncryptToBase64String(plain []byte) (string, error) {
	c, err := cs.Encrypt(plain)
	if err != nil {
		return "", app.CreateError(app.CryptographyError)
	}
	return base64.RawStdEncoding.EncodeToString(c), nil
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

func (cs CryptographyService) DecryptFromBase64String(base64Cipher string) (string, error) {
	c, err := base64.RawStdEncoding.DecodeString(base64Cipher)
	if err != nil {
		return "", app.CreateError(app.CryptographyError)
	}

	decrypted, err := cs.Decrypt(c)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
