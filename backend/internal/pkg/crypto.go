package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var backupEncryptionKey []byte

// InitBackupCrypto initializes the backup encryption key.
// It first checks the BACKUP_ENCRYPTION_KEY environment variable (hex-encoded),
// then tries to load from data/backup-encryption.key,
// and finally generates a new key and persists it if none exists.
func InitBackupCrypto(dataDir string) error {
	// 1. Check environment variable
	if envKey := os.Getenv("BACKUP_ENCRYPTION_KEY"); envKey != "" {
		key, err := hex.DecodeString(envKey)
		if err != nil {
			return fmt.Errorf("invalid BACKUP_ENCRYPTION_KEY: %w", err)
		}
		if len(key) != 32 {
			return fmt.Errorf("BACKUP_ENCRYPTION_KEY must be 32 bytes (64 hex chars)")
		}
		backupEncryptionKey = key
		return nil
	}

	// 2. Try to load from file
	keyFile := filepath.Join(dataDir, "backup-encryption.key")
	if data, err := os.ReadFile(keyFile); err == nil && len(data) > 0 {
		key, err := hex.DecodeString(string(data))
		if err != nil {
			return fmt.Errorf("invalid key file: %w", err)
		}
		if len(key) != 32 {
			return fmt.Errorf("key file must contain 32 bytes (64 hex chars)")
		}
		backupEncryptionKey = key
		log.Println("backup encryption key loaded from file")
		return nil
	}

	// 3. Generate and persist new key
	key, err := GenerateEncryptionKey()
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	if err := os.WriteFile(keyFile, []byte(hex.EncodeToString(key)), 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}
	backupEncryptionKey = key
	log.Println("backup encryption key generated and persisted to file")
	return nil
}

// GetBackupKey returns the current backup encryption key.
func GetBackupKey() []byte {
	return backupEncryptionKey
}

// GenerateEncryptionKey generates a random 32-byte key for AES-256.
func GenerateEncryptionKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt encrypts a plaintext string using AES-256-GCM with the given key.
// Returns base64(nonce[12] + ciphertext). The nonce is prepended to allow
// decryption without storing it separately.
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and seal
	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	// Prepend nonce to ciphertext and base64 encode
	result := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts a ciphertext produced by Encrypt using AES-256-GCM.
// The input must be base64(nonce[12] + encrypted_data).
func Decrypt(ciphertext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes")
	}

	// Base64 decode
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Split nonce and ciphertext
	nonce, cipherData := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}
