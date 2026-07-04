package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	block     cipher.Block
	initOnce  sync.Once
	initError error
)

func getBlock() (cipher.Block, error) {
	initOnce.Do(func() {
		keyHex := os.Getenv("ENCRYPTION_KEY")
		if keyHex == "" {
			initError = errors.New("ENCRYPTION_KEY environment variable is not set")
			log.Println("[crypto] WARNING: ENCRYPTION_KEY not set")
			return
		}

		keyBytes, err := hex.DecodeString(keyHex)
		if err != nil {
			initError = fmt.Errorf("ENCRYPTION_KEY must be a valid hex string: %v", err)
			log.Printf("[crypto] ERROR: %v", initError)
			return
		}
		if len(keyBytes) != 32 {
			initError = fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes (got %d)", len(keyBytes))
			log.Printf("[crypto] ERROR: %v", initError)
			return
		}

		b, err := aes.NewCipher(keyBytes)
		if err != nil {
			initError = fmt.Errorf("failed to create AES cipher: %v", err)
			log.Printf("[crypto] ERROR: %v", initError)
			return
		}
		block = b
	})

	if initError != nil {
		return nil, initError
	}
	if block == nil {
		return nil, errors.New("encryption not initialized")
	}
	return block, nil
}

// Encrypt encrypts plaintext using AES-256-GCM.
// Returns base64url-encoded nonce+ciphertext.
func Encrypt(plaintext string) (string, error) {
	b, err := getBlock()
	if err != nil {
		return "", fmt.Errorf("encryption error")
	}

	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return "", fmt.Errorf("encryption error")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encryption error")
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64url-encoded AES-256-GCM ciphertext.
func Decrypt(encoded string) (string, error) {
	b, err := getBlock()
	if err != nil {
		return "", fmt.Errorf("decryption error")
	}

	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decryption error")
	}

	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return "", fmt.Errorf("decryption error")
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("decryption error")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption error")
	}

	return string(plaintext), nil
}
