// Package encryptor provides AES-GCM encryption and decryption for
// sensitive environment variable values stored in .env files.
package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// Encryptor encrypts and decrypts environment variable values.
type Encryptor struct {
	key []byte
}

// New creates a new Encryptor from a passphrase. The passphrase is
// hashed with SHA-256 to produce a 32-byte AES key.
func New(passphrase string) *Encryptor {
	hash := sha256.Sum256([]byte(passphrase))
	return &Encryptor{key: hash[:]}
}

// Encrypt encrypts plaintext and returns a base64-encoded ciphertext.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded ciphertext and returns the plaintext.
func (e *Encryptor) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("encryptor: ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptMap encrypts all values in the provided map, returning a new map.
func (e *Encryptor) EncryptMap(env map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		enc, err := e.Encrypt(v)
		if err != nil {
			return nil, err
		}
		result[k] = enc
	}
	return result, nil
}

// DecryptMap decrypts all values in the provided map, returning a new map.
func (e *Encryptor) DecryptMap(env map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		dec, err := e.Decrypt(v)
		if err != nil {
			return nil, err
		}
		result[k] = dec
	}
	return result, nil
}
