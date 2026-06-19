package pkg

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// TestEncryptDecryptRoundTrip verifies that encrypting a plaintext and then
// decrypting the resulting ciphertext returns the original plaintext.
// ---------------------------------------------------------------------------
func TestEncryptDecryptRoundTrip(t *testing.T) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() error = %v", err)
	}

	plaintexts := []string{
		"hello world",
		"",
		"中文测试内容 🚀",
		strings.Repeat("a", 4096),
		"special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?",
	}

	for _, pt := range plaintexts {
		ciphertext, encErr := Encrypt(pt, key)
		if encErr != nil {
			t.Fatalf("Encrypt(%q) error = %v", pt, encErr)
		}
		if ciphertext == "" {
			t.Fatalf("Encrypt(%q) returned empty ciphertext", pt)
		}

		decrypted, decErr := Decrypt(ciphertext, key)
		if decErr != nil {
			t.Fatalf("Decrypt() error = %v", decErr)
		}
		if decrypted != pt {
			t.Fatalf("round-trip mismatch: got %q, want %q", decrypted, pt)
		}
	}
}

// ---------------------------------------------------------------------------
// TestDecryptInvalidCiphertext verifies that Decrypt correctly rejects
// invalid, tampered, or malformed ciphertexts.
// ---------------------------------------------------------------------------
func TestDecryptInvalidCiphertext(t *testing.T) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() error = %v", err)
	}

	// 1. Invalid base64.
	_, err = Decrypt("!!!not-valid-base64!!!", key)
	if err == nil {
		t.Error("expected error for invalid base64, got nil")
	}

	// 2. Valid base64 but tampered ciphertext.
	plaintext := "sensitive data"
	ciphertext, encErr := Encrypt(plaintext, key)
	if encErr != nil {
		t.Fatalf("Encrypt() error = %v", encErr)
	}

	raw, _ := base64.StdEncoding.DecodeString(ciphertext)
	// Flip a byte in the encrypted portion (after the nonce).
	if len(raw) > 13 {
		raw[13] ^= 0xFF
	}
	tampered := base64.StdEncoding.EncodeToString(raw)

	_, err = Decrypt(tampered, key)
	if err == nil {
		t.Error("expected error for tampered ciphertext, got nil")
	}

	// 3. Ciphertext too short (shorter than nonce size).
	_, err = Decrypt(base64.StdEncoding.EncodeToString([]byte{0x01}), key)
	if err == nil {
		t.Error("expected error for too-short ciphertext, got nil")
	}

	// 4. Wrong key size.
	_, err = Decrypt(ciphertext, []byte("short"))
	if err == nil {
		t.Error("expected error for wrong key size, got nil")
	}
}

// ---------------------------------------------------------------------------
// TestGenerateEncryptionKey verifies that the generated key is exactly 32
// bytes (AES-256) and that successive calls produce different keys.
// ---------------------------------------------------------------------------
func TestGenerateEncryptionKey(t *testing.T) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() error = %v", err)
	}
	if len(key) != 32 {
		t.Fatalf("key length = %d, want 32", len(key))
	}

	// Generate again; should be different (probabilistic).
	key2, err2 := GenerateEncryptionKey()
	if err2 != nil {
		t.Fatalf("second GenerateEncryptionKey() error = %v", err2)
	}
	if len(key2) != 32 {
		t.Fatalf("second key length = %d, want 32", len(key2))
	}

	// Extremely unlikely to collide for 32 random bytes.
	same := true
	for i := range key {
		if key[i] != key2[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("two generated keys are identical — suspicious")
	}
}

// ---------------------------------------------------------------------------
// TestInitBackupCrypto_EnvVar verifies that InitBackupCrypto correctly loads
// the encryption key from the BACKUP_ENCRYPTION_KEY environment variable.
// ---------------------------------------------------------------------------
func TestInitBackupCrypto_EnvVar(t *testing.T) {
	// Save the original key so we can restore it.
	origKey := backupEncryptionKey
	defer func() { backupEncryptionKey = origKey }()

	// Generate a valid 32-byte key and hex-encode it (64 hex chars).
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() error = %v", err)
	}
	hexKey := hex.EncodeToString(key)

	// Set the environment variable.
	_ = os.Setenv("BACKUP_ENCRYPTION_KEY", hexKey)
	defer func() { _ = _ = os.Unsetenv("BACKUP_ENCRYPTION_KEY") }()

	// Use a temp directory that doesn't contain a key file.
	tmpDir := t.TempDir()

	// Call InitBackupCrypto — it should pick up the env var.
	if err := InitBackupCrypto(tmpDir); err != nil {
		t.Fatalf("InitBackupCrypto() error = %v", err)
	}

	// Verify the loaded key matches.
	loaded := GetBackupKey()
	if len(loaded) != 32 {
		t.Fatalf("loaded key length = %d, want 32", len(loaded))
	}
	for i := range key {
		if loaded[i] != key[i] {
			t.Fatalf("loaded key mismatch at byte %d: got %02x, want %02x", i, loaded[i], key[i])
		}
	}

	// Also test that an invalid hex env var returns an error.
	origKey2 := backupEncryptionKey
	backupEncryptionKey = nil
	_ = os.Setenv("BACKUP_ENCRYPTION_KEY", "not-valid-hex")
	err = InitBackupCrypto(tmpDir)
	if err == nil {
		t.Error("expected error for invalid hex env var, got nil")
	}
	backupEncryptionKey = origKey2
}

// ---------------------------------------------------------------------------
// TestInitBackupCrypto_File verifies InitBackupCrypto loads from file when
// the env var is not set.
// ---------------------------------------------------------------------------
func TestInitBackupCrypto_File(t *testing.T) {
	origKey := backupEncryptionKey
	defer func() { backupEncryptionKey = origKey }()

	// Ensure env var is not set.
	_ = os.Unsetenv("BACKUP_ENCRYPTION_KEY")

	tmpDir := t.TempDir()

	// Write a valid key file.
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() error = %v", err)
	}
	hexKey := hex.EncodeToString(key)
	keyFile := filepath.Join(tmpDir, "backup-encryption.key")
	if err := os.WriteFile(keyFile, []byte(hexKey), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := InitBackupCrypto(tmpDir); err != nil {
		t.Fatalf("InitBackupCrypto() error = %v", err)
	}

	loaded := GetBackupKey()
	if len(loaded) != 32 {
		t.Fatalf("loaded key length = %d, want 32", len(loaded))
	}
}

// ---------------------------------------------------------------------------
// TestInitBackupCrypto_Generate verifies InitBackupCrypto generates and
// persists a new key when neither env var nor file exists.
// ---------------------------------------------------------------------------
func TestInitBackupCrypto_Generate(t *testing.T) {
	origKey := backupEncryptionKey
	defer func() { backupEncryptionKey = origKey }()

	_ = os.Unsetenv("BACKUP_ENCRYPTION_KEY")

	tmpDir := t.TempDir()

	if err := InitBackupCrypto(tmpDir); err != nil {
		t.Fatalf("InitBackupCrypto() error = %v", err)
	}

	loaded := GetBackupKey()
	if len(loaded) != 32 {
		t.Fatalf("loaded key length = %d, want 32", len(loaded))
	}

	// Verify the key file was created.
	keyFile := filepath.Join(tmpDir, "backup-encryption.key")
	data, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("key file not created: %v", err)
	}
	if len(data) != 64 {
		t.Fatalf("key file hex length = %d, want 64", len(data))
	}
}
