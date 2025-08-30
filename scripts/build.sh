#!/bin/bash
set -e

echo "Building Super MCP Server..."

# Check if Super binary exists
if [ ! -f "binaries/super-linux-amd64" ]; then
    echo "Error: Super binary not found. Run './scripts/setup.sh' first."
    exit 1
fi

# Clean previous builds
rm -f super-mcp-server bootstrap lambda-deployment.zip

# Build for local development
echo "Building for local development..."
go mod tidy
go build -o super-mcp-server main.go
echo "✓ Local binary created: super-mcp-server"

# Build for Lambda
echo "Building for AWS Lambda..."
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip lambda-deployment.zip bootstrap
rm bootstrap
echo "✓ Lambda deployment package created: lambda-deployment.zip"

# Show file sizes
echo
echo "Build Summary:"
echo "Local binary size: $(du -h super-mcp-server | cut -f1)"
echo "Lambda package size: $(du -h lambda-deployment.zip | cut -f1)"
echo
echo "Ready to deploy!"
echo "Local testing: ./super-mcp-server"
echo "Lambda deployment: Upload lambda-deployment.zip to AWS Lambda"