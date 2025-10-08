// Copyright (c) 2025 William Chastain
// Licensed under the MIT License. See LICENSE.txt file in the project root for details.

package auth

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// PromptForPassword securely prompts the user for a password without echoing input.
// If confirm is true, requires password confirmation.
func PromptForPassword(confirm bool) (string, error) {
	fmt.Print("::Enter password: ")
	pwdBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	} else if len(pwdBytes) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}
	fmt.Println() // Move to the next line after password input

	if confirm {
		fmt.Print("::Confirm password: ")
		confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println() // Move to the next line after confirmation input

		if string(pwdBytes) != string(confirmBytes) {
			return "", fmt.Errorf("passwords do not match")
		}
	}

	return string(pwdBytes), nil
}
