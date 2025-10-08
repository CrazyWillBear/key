// Copyright (c) 2025 William Chastain
// Licensed under the MIT License. See LICENSE.txt file in the project root for details.

package key

import (
	"fmt"
	"key/auth"
	"key/cryptoutil"
	"os"
)

type Key struct {
	Path              string // path to key file
	EncryptedKeyBytes []byte // ciphertext
	DecryptedKeyBytes []byte // plaintext
}

// Clear wipes the decrypted key from memory
func (k *Key) Clear() {
	// Wipe the decrypted key from memory
	for i := range k.DecryptedKeyBytes {
		k.DecryptedKeyBytes[i] = 0
	}
	k.DecryptedKeyBytes = nil
}

// Decrypt decrypts the key using the provided password
func (k *Key) Decrypt(password string) error {
	// Decrypt the key
	decryptedKey, err := cryptoutil.DecryptWithPassword(k.EncryptedKeyBytes, password)
	if err != nil {
		return fmt.Errorf("failed to decrypt key: %w", err)
	}

	k.DecryptedKeyBytes = decryptedKey
	return nil
}

// CreateNewKey creates a new encrypted key and saves it to the specified path
func CreateNewKey(keyPath string) error {
	// Get password from user
	fmt.Println("::Choose a password to encrypt your key. IT IS VERY IMPORTANT YOU REMEMBER THIS PASSWORD!")
	password, err := auth.PromptForPassword(true)
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Create new encrypted key
	k, err := newEncryptedKey(keyPath, password)
	if err != nil {
		return fmt.Errorf("failed to create new encrypted key: %w", err)
	}

	// Write encrypted key to file
	err = os.WriteFile(k.Path, k.EncryptedKeyBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write key to file: %w", err)
	}

	return nil
}

// LoadKey loads a Key from the specified path
func LoadKey(path string) *Key {
	// Read the key file
	data, err := os.ReadFile(path)
	if err != nil {
		// Return empty Key if file doesn't exist or can't be read
		return nil
	}

	// Make Key
	k := &Key{Path: path, EncryptedKeyBytes: data, DecryptedKeyBytes: nil}

	return k
}

// newEncryptedKey creates a new Key, encrypts it with the provided password, and returns it
func newEncryptedKey(path string, password string) (*Key, error) {
	// Generate random key
	keyBytes, err := cryptoutil.GenerateRandomKey(32) // 256-bit key
	if err != nil {
		return nil, err
	}

	// Encrypt key with password
	keyBytes, err = cryptoutil.EncryptWithPassword(keyBytes, password)
	if err != nil {
		return nil, err
	}

	// Make Key
	k := &Key{Path: path, EncryptedKeyBytes: keyBytes, DecryptedKeyBytes: nil}

	return k, nil
}
