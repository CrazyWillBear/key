#!/bin/bash

# Build the Go binary with necessary flags
echo "Building Key..."
go build -ldflags="-s -w" -o "bin/key" .

# Check if build was successful
if [ ! -f "bin/key" ]; then
    echo "Build failed!"
    exit 1
fi

# Get the current user's home directory
USER_HOME=$(eval echo ~"$USER")
KEY_DIR="$USER_HOME/.key"

# Create the .key directory
echo "Creating $KEY_DIR directory..."
mkdir -p "$KEY_DIR"

# Move the binary to /usr/bin (requires sudo)
echo "Installing Key to /usr/bin/..."
sudo mv "bin/key" "/usr/bin/key"

# Make it executable
sudo chmod +x "/usr/bin/key"

echo "Installation complete!"
