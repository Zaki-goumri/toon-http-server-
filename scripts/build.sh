#!/bin/bash

set -e

echo "========================================"
echo "Building TOON vs JSON Servers"
echo "========================================"
echo ""

# Create bin directory if it doesn't exist
mkdir -p bin

# Build TOON server
echo "Building TOON server..."
go build -o bin/toon-server cmd/toon-server/main.go
echo "✓ Built bin/toon-server"

# Build JSON server
echo "Building JSON server..."
go build -o bin/json-server cmd/json-server/main.go
echo "✓ Built bin/json-server"

echo ""
echo "========================================"
echo "Build Complete!"
echo "========================================"
echo ""
echo "Run servers:"
echo "  ./bin/toon-server  (port 8080)"
echo "  ./bin/json-server  (port 8081)"
