// Copyright (c) 2025 William Chastain
// Licensed under the MIT License. See LICENSE.txt file in the project root for details.

package vault

import (
	"fmt"
	"key/cryptoutil"
	"os"
)

// Vault represents a password raw_vault with its path, content, and logins
type Vault struct {
	Path             string
	EncryptedContent []byte
	DecryptedContent []byte
}

// LoadVault creates a new Vault instance with the given path
func LoadVault(path string) (*Vault, error) {
	// Read the encrypted content from the file
	rawContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Check if the content is encrypted by looking for the magic header
	encrypted := len(rawContent) >= 8 && string(rawContent[:8]) == "KEYVLT01"

	if !encrypted {
		// If the vault is not encrypted, set DecryptedContent to rawContent
		return &Vault{Path: path, EncryptedContent: nil, DecryptedContent: rawContent}, nil
	} else {
		// If the vault is encrypted, set EncryptedContent to rawContent
		return &Vault{Path: path, EncryptedContent: rawContent, DecryptedContent: nil}, nil
	}
}

// Unlock decrypts the raw_vault content using the provided key and populates the Logins map
func (v *Vault) Unlock(key []byte) error {
	if v.DecryptedContent != nil {
		// Already decrypted
		return fmt.Errorf("vault is already unlocked")
	}

	// Decrypt the content
	err := cryptoutil.DecryptFileWithKey(v.Path, v.Path, key)
	if err != nil {
		return err
	}

	return nil
}

// Lock encrypts the Logins map using the provided key and writes it back to the raw_vault file
func (v *Vault) Lock(key []byte) error {
	if v.EncryptedContent != nil {
		// Already encrypted
		return fmt.Errorf("vault is already locked")
	}

	err := cryptoutil.EncryptFileWithKey(v.Path, v.Path, key)
	if err != nil {
		return err
	}

	return nil
}
