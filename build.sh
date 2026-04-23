#!/bin/bash

# Build script for AI Coder Context Server
# Cross-platform builds for Windows, macOS, and Linux

set -e

VERSION=${VERSION:-"1.0.0"}
OUTPUT_DIR="dist"

echo "Building AI Coder Context Server v${VERSION}"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build function
build() {
    GOOS=$1 GOARCH=$2 EXT=$3
    OUTPUT="$OUTPUT_DIR/ai-coder-context-server-${VERSION}-${GOOS}-${GOARCH}"
    if [ "$EXT" != "" ]; then
        OUTPUT="$OUTPUT.$EXT"
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUTPUT" ./cmd/server
}

# Build for all platforms
build "windows" "amd64" "exe"
build "darwin" "amd64" ""
build "darwin" "arm64" ""
build "linux" "amd64" ""
build "linux" "arm64" ""

echo ""
echo "Build complete! Output in $OUTPUT_DIR:"
ls -la "$OUTPUT_DIR"