// Key - a command-line utility for securely encrypting and decrypting files
// Copyright (c) 2025 William Chastain
//
// Author: William Chastain (github.com/CrazyWillBear - www.capbear.net - williamchastain2005@gmail.com)
// License: MIT License (see LICENSE.txt for full text)
// Homepage: https://github.com/CrazyWillBear/key
//
// Description:
//   A command-line utility for securely encrypting and decrypting files using a key + password combo. It uses AES-GCM
//   for encryption and PBKDF2 for key derivation. It supports creating new keys, locking files into vaults, and
//   unlocking them.

package main

import (
	"fmt"
	"key/commands"
	"key/config"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return
	}

	// Execute commands
	if err := commands.Execute(*cfg); err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
}
