#!/bin/bash

set -e

# Release script for tusk CLI
# Builds binaries for multiple platforms

VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X github.com/sboy99/projects.go/tusk/internal/version.BuildVersion=${VERSION}"
LDFLAGS="${LDFLAGS} -X github.com/sboy99/projects.go/tusk/internal/version.BuildCommit=${COMMIT}"
LDFLAGS="${LDFLAGS} -X github.com/sboy99/projects.go/tusk/internal/version.BuildTime=${BUILD_TIME}"

# Platforms to build for
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

mkdir -p dist

for PLATFORM in "${PLATFORMS[@]}"; do
    PLATFORM_SPLIT=(${PLATFORM//\// })
    GOOS=${PLATFORM_SPLIT[0]}
    GOARCH=${PLATFORM_SPLIT[1]}
    
    OUTPUT_NAME="tusk"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="tusk.exe"
    fi
    
    OUTPUT_PATH="dist/tusk-${VERSION}-${GOOS}-${GOARCH}/${OUTPUT_NAME}"
    
    echo "Building for ${GOOS}/${GOARCH}..."
    GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o "${OUTPUT_PATH}" ./cmd/tusk
    
    echo "Built: ${OUTPUT_PATH}"
done

echo "Release build complete! Binaries are in dist/"

