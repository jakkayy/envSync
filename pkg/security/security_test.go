package security

import (
	"testing"
)

func TestSecretDetector(t *testing.T) {
	envMap := map[string]string{
		"AWS_SECRET": "AKIAIOSFODNN7EXAMPLE",
		"DB_URL":     "postgres://user:password123@localhost:5432/mydb",
		"NORMAL_VAR": "hello_world",
	}

	findings := DetectSecrets(envMap)
	if len(findings) != 2 {
		t.Fatalf("expected 2 secret findings, got %d", len(findings))
	}
}

func TestSecretMasker(t *testing.T) {
	maskedPass := MaskValue("DATABASE_PASSWORD", "super_secret_password_123")
	if maskedPass == "super_secret_password_123" {
		t.Errorf("password was not masked")
	}

	normal := MaskValue("PORT", "8080")
	if normal != "8080" {
		t.Errorf("normal value should not be masked, got %s", normal)
	}
}
