#!/bin/bash
set -e

echo "Building Super from source for Linux..."
mkdir -p binaries temp

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is required to build Super. Please install Go first."
    exit 1
fi

# Clone and build Super for Linux
echo "Cloning Super repository..."
cd temp
git clone https://github.com/brimdata/super.git
cd super

echo "Building Super for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o ../../binaries/super-linux-amd64 ./cmd/super

# Clean up
cd ../..
rm -rf temp

# Make executable
chmod +x binaries/super-linux-amd64

echo "Super binary built successfully at binaries/super-linux-amd64"
echo "File size: $(du -h binaries/super-linux-amd64 | cut -f1)"
echo "Setup complete! You can now build the project."