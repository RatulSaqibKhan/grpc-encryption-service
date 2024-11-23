package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"os"
)

// DecryptData decrypts the given data using AES-256-CTR algorithm.
func Decrypt(data string) (string, error) {
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

	// Decode the encrypted data from hexadecimal
	encryptedText, err := hex.DecodeString(data)
	if err != nil {
		return "", errors.New("failed to decode encrypted text")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(password)
	if err != nil {
		return "", errors.New("failed to create AES cipher block")
	}

	// Create a stream cipher for decryption
	stream := cipher.NewCTR(block, iv)

	// Decrypt the data
	decrypted := make([]byte, len(encryptedText))
	stream.XORKeyStream(decrypted, encryptedText)

	return string(decrypted), nil
}
