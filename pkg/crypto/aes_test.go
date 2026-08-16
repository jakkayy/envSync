package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptAES256GCM(t *testing.T) {
	passphrase := "my-super-secret-master-password"
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}

	key := DeriveKey(passphrase, salt)
	if len(key) != 32 {
		t.Fatalf("expected key length 32 bytes, got %d", len(key))
	}

	plaintext := []byte("DATABASE_URL=postgres://user:pass@localhost:5432/db")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatalf("ciphertext should not match plaintext")
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted text %s does not match original %s", string(decrypted), string(plaintext))
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := DeriveKey("pass1", []byte("1234567890123456"))
	key2 := DeriveKey("pass2", []byte("1234567890123456"))

	plaintext := []byte("SECRET_API_KEY=xyz123")
	ciphertext, err := Encrypt(plaintext, key1)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(ciphertext, key2)
	if err == nil {
		t.Errorf("expected error when decrypting with wrong key, got nil")
	}
}
