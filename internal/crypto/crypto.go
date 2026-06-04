package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	pbkdf2Iterations = 100_000
	keySize          = 32 // AES-256
	nonceSize        = 12 // GCM standard nonce
	saltSize         = 32 // 256-bit salt
)

// GenerateSalt returns a random hex-encoded 32-byte salt.
func GenerateSalt() (string, error) {
	b := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// DeriveKey derives a 32-byte AES key from password and hex-encoded salt using PBKDF2-SHA256.
// The resulting key is held in memory only; never persisted to disk.
func DeriveKey(password, saltHex string) ([]byte, error) {
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return nil, err
	}
	return pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, keySize, sha256.New), nil
}

// Encrypt encrypts plaintext with AES-256-GCM.
// Returns base64(nonce || ciphertext). A fresh nonce is generated per call.
func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, nonceSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	// Seal appends ciphertext (including GCM auth tag) after the nonce.
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt decrypts a base64-encoded blob produced by Encrypt.
func Decrypt(key []byte, encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
