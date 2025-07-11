package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Cryptographic constants.
const (
	// NonceSize is the size of the nonce used for AES-GCM encryption.
	NonceSize = 12
)

func AesEncrypt(key []byte, plaintext string) ([]byte, []byte, error) {
	block, errBlock := aes.NewCipher(key)
	if errBlock != nil {
		return nil, nil, errBlock
	}

	nonce := make([]byte, NonceSize)

	if _, errRand := io.ReadFull(rand.Reader, nonce); errRand != nil {
		return nil, nonce, errRand
	}

	aesgcm, errGCM := cipher.NewGCM(block)
	if errGCM != nil {
		return nil, nonce, errGCM
	}

	return aesgcm.Seal(nil, nonce, []byte(plaintext), nil), nonce, nil
}

func AesDecrypt(key []byte, nonce []byte, ciphertext []byte) (string, error) {
	block, errBlock := aes.NewCipher(key)
	if errBlock != nil {
		return "", errBlock
	}

	aesgcm, errGCM := cipher.NewGCM(block)
	if errGCM != nil {
		return "", errGCM
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	return string(plaintext), err
}
