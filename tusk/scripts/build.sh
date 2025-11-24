#!/bin/bash

set -e

# Build script for tusk CLI

VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X github.com/sboy99/projects.go/tusk/internal/version.BuildVersion=${VERSION}"
LDFLAGS="${LDFLAGS} -X github.com/sboy99/projects.go/tusk/internal/version.BuildCommit=${COMMIT}"
LDFLAGS="${LDFLAGS} -X github.com/sboy99/projects.go/tusk/internal/version.BuildTime=${BUILD_TIME}"

# Build for current platform
echo "Building tusk..."
go build -ldflags "${LDFLAGS}" -o bin/tusk ./cmd/tusk

echo "Build complete: bin/tusk"

