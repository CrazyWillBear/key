<!-- TOC -->
* [Key](#key)
  * [Overview](#overview)
    * [Summary](#summary)
    * [TL;DR](#tldr)
  * [Installation](#installation)
  * [Usage](#usage)
    * [Generate a New Key](#generate-a-new-key)
    * [Lock (Encrypt) a File](#lock-encrypt-a-file)
    * [Unlock (Decrypt) a File](#unlock-decrypt-a-file)
    * [Use Verbose Mode](#use-verbose-mode)
    * [Use a Custom Key Path](#use-a-custom-key-path)
  * [Configuration](#configuration)
  * [Security](#security)
  * [Development](#development)
    * [Build](#build)
  * [License](#license)
<!-- TOC -->

# Key

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a6b66a290bab42038701fb2cea8ab0f6)](https://app.codacy.com/gh/CrazyWillBear/key?utm_source=github.com&utm_medium=referral&utm_content=CrazyWillBear/key&utm_campaign=Badge_Grade)

A command-line utility for securely encrypting and decrypting files using a key + password combo.

## Overview

### Summary

Key is a simple command-line tool for securely encrypting and decrypting files. It uses both a personal key and a
password to encrypt and decrypt files, ensuring that your data remains safe even if your key file is compromised.
Through this two-factor approach, Key provides robust security for your sensitive files.

Furthermore, Key is designed to be user-friendly, with clear commands such as `lock` and `unlock` to manage file
encryption and decryption. It also supports configuration via a config file, environment variables, or command-line
flags, allowing for flexible usage.

Key employs strong encryption standards, including AES-GCM for file encryption and PBKDF2 for password-based key
derivation. It ensures that keys are securely wiped from memory after use and that passwords are *never* stored.

### TL;DR

Key is a file encryption tool that uses a two-factor approach:
1. A personal encryption key
2. A password protecting that key

This approach ensures that even if your key file is compromised, your data remains secure without the password. Key is
easy to use, with simple commands for locking and unlocking files, and it supports configuration through various means.

## Installation

```sh
# Build from source
go build -o bin/key

# Add to your PATH (optional)
cp bin/key /usr/local/bin/
```

## Usage

### Generate a New Key

Before first use, generate your personal encryption key:

```sh
key newkey
```

This creates an encrypted key file at the configured location (default: `~/.key/key.pem`).

### Lock (Encrypt) a File

```sh
key lock myfile.txt
```

### Unlock (Decrypt) a File

```sh
key unlock myfile.txt
```

### Use Verbose Mode

Add the `-v` flag to see detailed operation steps:

```sh
key -v unlock myfile.txt
```

### Use a Custom Key Path

1. Key allows specifying a custom key file path with the `--key-path` flag:

```sh
key --key-path /path/to/my/key.pem unlock myfile.txt
```

2. Or with the `KEY_KEY_PATH` environment variable:

```sh
export KEY_KEY_PATH=/path/to/my/key.pem
```

3. Or by setting it in the configuration file (see below).

## Configuration

Key uses a configuration file located at `~/.key/config.toml`. It simply contains the path to your key file.

```toml
# Key Manager Configuration
key_path = "$HOME/.key/key.pem"
```

## Security

- Keys are encrypted using AES-GCM with password-derived keys (PBKDF2)
- Files are encrypted with AES-GCM using the decrypted key
- Keys are wiped from memory after use
- Passwords are *never* stored

## Development

- Built with Go 1.25+
- Uses Cobra for CLI commands
- Uses Viper for configuration
- Dependencies managed with Go modules

### Build

```sh
go build -o bin/key
```

## License

This project is licensed under the [MIT License](LICENSE.txt).