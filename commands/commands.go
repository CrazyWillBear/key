// Copyright (c) 2025 William Chastain
// Licensed under the MIT License. See LICENSE.txt file in the project root for details.

package commands

import (
	"key/auth"
	"key/config"
	"key/key"
	"key/vault"

	"github.com/spf13/cobra"
)

// Global configuration variable
var globalConfig config.Config

// Flags
var keyPathOverride string
var verbose bool

// Root command
var rootCmd = &cobra.Command{
	Use:   "key",
	Short: "Use a 'key' to lock, unlock, and create vaults",
}

// Command to unlock the vault
var unlockCmd = &cobra.Command{
	Use:   "unlock [file]",
	Short: "Unlock a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get password
		password, err := auth.PromptForPassword(false)
		if err != nil {
			cmd.PrintErrln("::Error getting password:", err)
			return
		}

		cmd.Println("::Unlocking file...")

		// Load vault
		if verbose {
			cmd.Print("\t- Loading vault...")
		}
		usrVault, err := vault.LoadVault(args[0])
		if err != nil {
			cmd.PrintErrln("\r::Error loading vault:", err)
			return
		}
		if verbose {
			cmd.Println("\r\t- Vault loaded.   ")
		}

		// Load key
		if verbose {
			cmd.Print("\t- Loading key...")
		}
		k := key.LoadKey(globalConfig.KeyPath)
		if k == nil {
			cmd.PrintErrln("\r::Key cannot be found, have you run `key newkey`?.")
			return
		}
		if verbose {
			cmd.Println("\r\t- Key loaded.   ")
		}

		// Decrypt key
		err = k.Decrypt(password)
		if err != nil {
			cmd.PrintErrln("::Error decrypting key:", err)
			return
		}
		if verbose {
			cmd.Println("\t- Key decrypted.")
		}

		// Unlock vault with key
		if verbose {
			cmd.Print("\t- Unlocking vault with key...")
		}
		err = usrVault.Unlock(k.DecryptedKeyBytes)
		if err != nil {
			cmd.PrintErrln("\r::Error unlocking vault:       ", err)
			return
		}
		if verbose {
			cmd.Println("\r\t- Vault unlocked.            ")
		}

		// Clear key from memory
		if verbose {
			cmd.Print("\t- Clearing key from memory...")
		}
		k.Clear()
		if verbose {
			cmd.Println("\r\t- Key cleared.               ")
		}

		println("::File unlocked successfully!")
	},
}

var lockCmd = &cobra.Command{
	Use:   "lock [file]",
	Short: "Lock a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get password
		password, err := auth.PromptForPassword(false)
		if err != nil {
			cmd.PrintErrln("::Error getting password:", err)
			return
		}

		cmd.Println("::Locking file...")

		// Create vault
		if verbose {
			cmd.Print("\t- Creating vault...")
		}
		usrVault, err := vault.LoadVault(args[0])
		if err != nil {
			cmd.PrintErrln("\r::Error loading vault:", err)
			return
		}
		if verbose {
			cmd.Println("\r\t- Vault created.   ")
		}

		// Load key
		if verbose {
			cmd.Print("\t- Loading key...")
		}
		k := key.LoadKey(globalConfig.KeyPath)
		if k == nil {
			cmd.PrintErrln("\r::Error loading key, have you ran `key newkey`?")
			return
		}
		if verbose {
			cmd.Println("\r\t- Key loaded.   ")
		}

		// Decrypt key
		err = k.Decrypt(password)
		if err != nil {
			cmd.PrintErrln("::Error decrypting key:", err)
			return
		}
		if verbose {
			cmd.Println("\t- Key decrypted.")
		}

		// Lock vault with key
		if verbose {
			cmd.Print("\t- Locking vault with key...")
		}
		err = usrVault.Lock(k.DecryptedKeyBytes)
		if err != nil {
			cmd.PrintErrln("\r::Error locking vault:       ", err)
			return
		}
		if verbose {
			cmd.Println("\r\t- Vault locked.            ")
		}

		// Clear key from memory
		if verbose {
			cmd.Print("\t- Clearing key from memory...")
		}
		k.Clear()
		if verbose {
			cmd.Println("\r\t- Key cleared.               ")
		}

		cmd.Println("::File locked successfully.")
	},
}

// Command to generate a new key file
var newKeyCmd = &cobra.Command{
	Use:   "newkey",
	Short: "Generate a new key at your configured key path",
	Run: func(cmd *cobra.Command, args []string) {
		// Check for existing key
		k := key.LoadKey(globalConfig.KeyPath)
		if k != nil {
			cmd.Println("::WARNING!!! Generating a new key will overwrite your existing key, and you will lose access to any vaults encrypted with the old key.")
			cmd.Println("::Enter your existing key password to continue, or Ctrl+C to abort.")

			// Decrypt existing key to verify password
			password, err := auth.PromptForPassword(false)
			if err != nil {
				cmd.PrintErrln("\t- Error getting password:", err)
				return
			}

			err = k.Decrypt(password)
			if err != nil {
				cmd.PrintErrln("\t- Incorrect password, aborting...", err)
				return
			}

			cmd.Println("\t- Existing key verified, continuing...")
		}

		// Create a new key and save it to the specified path
		err := key.CreateNewKey(globalConfig.KeyPath)
		if err != nil {
			cmd.Println("Error generating new key:", err)
			return
		}
		cmd.Println("::New key generated and saved to", globalConfig.KeyPath)
	},
}

// Initialize commands and flags
func init() {
	// Flags
	rootCmd.PersistentFlags().StringVar(&keyPathOverride, "key-path", "", "Override key file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Commands
	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
	rootCmd.AddCommand(newKeyCmd)
}

// Execute runs the root command with the provided configuration
func Execute(cfg config.Config) error {
	globalConfig = cfg

	if keyPathOverride != "" {
		globalConfig.KeyPath = keyPathOverride
	}

	return rootCmd.Execute()
}
