package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize   = 16
	KeyLength  = 32
	Iterations = 100000
)

// GenerateSalt creates a 16-byte random salt for PBKDF2
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %w", err)
	}
	return salt, nil
}

// DeriveKey derives a 32-byte AES key from a passphrase and salt using PBKDF2-HMAC-SHA256
func DeriveKey(passphrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, Iterations, KeyLength, sha256.New)
}
