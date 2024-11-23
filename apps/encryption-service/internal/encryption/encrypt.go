package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"os"
)

// Encrypt encrypts the given plaintext using AES-256-CTR algorithm.
func Encrypt(plaintext string) (string, error) {
	// Load password and IV from environment variables
	password := []byte(os.Getenv("ENCRYPTION_PASSWORD"))
	iv := []byte(os.Getenv("ENCRYPTION_IV"))
	algorithm := os.Getenv("ENCRYPTION_ALGORITHM") // Should be "aes-256-ctr"

	// Validate input
	if len(password) != 32 {
		return "", errors.New("invalid encryption key size: must be 32 bytes")
	}
	if len(iv) != aes.BlockSize {
		return "", errors.New("invalid IV size: must be 16 bytes")
	}
	if algorithm != "aes-256-ctr" {
		return "", errors.New("unsupported encryption algorithm: must be aes-256-ctr")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(password)
	if err != nil {
		return "", errors.New("failed to create AES cipher block")
	}

	// Create a stream cipher for encryption
	stream := cipher.NewCTR(block, iv)

	// Encrypt the plaintext
	encrypted := make([]byte, len(plaintext))
	stream.XORKeyStream(encrypted, []byte(plaintext))

	// Return the encrypted data as a hexadecimal string
	return hex.EncodeToString(encrypted), nil
}
