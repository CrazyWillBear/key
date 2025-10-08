// Copyright (c) 2025 William Chastain
// Licensed under the MIT License. See LICENSE.txt file in the project root for details.

package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptWithPassword encrypts plaintext using a password-derived key with AES-GCM.
// The salt and nonce are prepended to the ciphertext.
func EncryptWithPassword(plaintext []byte, password string) ([]byte, error) {
	// Generate a random salt for key derivation
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive key from password using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)

	// Create AES cipher and GCM mode
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt and prepend salt and nonce to ciphertext
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	result := make([]byte, 16+len(nonce)+len(ciphertext))
	copy(result, salt)
	copy(result[16:], nonce)
	copy(result[16+len(nonce):], ciphertext)

	return result, nil
}

// DecryptWithPassword decrypts ciphertext using a password-derived key with AES-GCM.
// Expects salt and nonce to be prepended to the ciphertext.
func DecryptWithPassword(data []byte, password string) ([]byte, error) {
	if len(data) < 28 { // 16 bytes salt + 12 bytes nonce minimum
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract salt, nonce, and ciphertext
	salt := data[:16]
	nonce := data[16:28]
	ciphertext := data[28:]

	// Derive key from password using same parameters
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)

	// Create AES cipher and GCM mode
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt and handle the error
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: wrong password or corrupted data")
	}
	return plaintext, nil
}

// Encrypt encrypts the plaintext using AES-GCM and returns the ciphertext with a prepended nonce.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	// Create a new AES cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Create a GCM (Galois/Counter Mode) cipher from the AES block
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// Generate a random nonce of the correct size for GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// Prepend nonce to ciphertext
	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts the ciphertext using AES-GCM and returns the plaintext.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	// Create a new AES cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Create a GCM cipher from the AES block
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// Extract the nonce from the beginning of the ciphertext
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	// Decrypt the ciphertext using the nonce and return the plaintext
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// GenerateRandomKey generates a random key of the specified size in bytes.
func GenerateRandomKey(size int) ([]byte, error) {
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
