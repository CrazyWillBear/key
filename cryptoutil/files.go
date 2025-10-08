// Copyright (c) 2025 William Chastain
// Licensed under the MIT License. See LICENSE.txt file in the project root for details.

package cryptoutil

import (
	"os"
)

// EncryptFileWithKey encrypts a file using a provided key and adds a magic header
func EncryptFileWithKey(inputPath, outputPath string, key []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		return err
	}

	// Prepend magic header for file identification
	magicHeader := []byte("KEYVLT01")
	result := make([]byte, 8+len(ciphertext))
	copy(result, magicHeader)
	copy(result[8:], ciphertext)

	return os.WriteFile(outputPath, result, 0644)
}

// DecryptFileWithKey decrypts a file using a provided key and handles magic header detection
func DecryptFileWithKey(inputPath, outputPath string, key []byte) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Check for and remove magic header if present
	if len(data) >= 8 && string(data[:8]) == "KEYVLT01" {
		data = data[8:] // Remove header
	}

	plaintext, err := Decrypt(data, key)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, plaintext, 0644)
}
