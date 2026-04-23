# Codebase Brain - Build Guide

Instructions for building Codebase Brain on various platforms.

---

## Prerequisites

- **Go 1.23+** - Download from https://go.dev/dl/

Verify installation:
```bash
go version
# Should show: go1.23.x
```

---

## Quick Build

```bash
# Clone the repository
git clone https://github.com/your-repo/codebase-brain.git
cd codebase-brain

# Download dependencies
go mod download

# Build for current platform
go build -o codebase-brain ./cmd/server

# Run
./codebase-brain
```

---

## Build Targets

Using the included Makefile:

| Target | Description |
|--------|-------------|
| `make build` | Build for current OS |
| `make all` | Build for all platforms |
| `make clean` | Remove build artifacts |
| `make test` | Run tests |
| `make install` | Install to $GOPATH/bin |

---

## Manual Builds

### Windows

```powershell
# AMD64 (x86-64)
GOOS=windows GOARCH=amd64 go build -o codebase-brain.exe ./cmd/server

# ARM64
GOOS=windows GOARCH=arm64 go build -o codebase-brain-arm64.exe ./cmd/server
```

### macOS

```bash
# AMD64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o codebase-brain-darwin-amd64 ./cmd/server

# ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o codebase-brain-darwin-arm64 ./cmd/server
```

### Linux

```bash
# AMD64
GOOS=linux GOARCH=amd64 go build -o codebase-brain-linux-amd64 ./cmd/server

# ARM64
GOOS=linux GOARCH=arm64 go build -o codebase-brain-linux-arm64 ./cmd/server
```

---

## Build Scripts

### Windows (PowerShell)

```powershell
./build.ps1
```

This produces:
- `dist/codebase-brain-1.0.0-windows-amd64.exe`

### macOS/Linux (Bash)

```bash
./build.sh
```

This produces:
- `dist/codebase-brain-1.0.0-darwin-amd64`
- `dist/codebase-brain-1.0.0-darwin-arm64`
- `dist/codebase-brain-1.0.0-linux-amd64`
- `dist/codebase-brain-1.0.0-linux-arm64`

---

## Cross-Compilation

Build for all platforms from one machine:

```bash
# Using build.sh
./build.sh

# Or manually:
GOOS=darwin GOARCH=amd64 go build -o dist/codebase-brain-darwin-amd64 ./cmd/server
GOOS=darwin GOARCH=arm64 go build -o dist/codebase-brain-darwin-arm64 ./cmd/server
GOOS=linux GOARCH=amd64 go build -o dist/codebase-brain-linux-amd64 ./cmd/server
GOOS=linux GOARCH=arm64 go build -o dist/codebase-brain-linux-arm64 ./cmd/server
GOOS=windows GOARCH=amd64 go build -o dist/codebase-brain-windows-amd64.exe ./cmd/server
```

---

## Installation

### Local Installation

```bash
# Move to PATH
cp codebase-brain /usr/local/bin/

# Or install via make
make install
```

### Verification

```bash
codebase-brain --help
```

---

## Build Options

### Build with Tags

```bash
# Enable debug symbols
go build -tags debug -o codebase-brain ./cmd/server

# Use different config
go build -o codebase-brain -ldflags="-s -w" ./cmd/server
```

### Version Info

The version is set in `cmd/server/main.go`:
```go
s := server.NewMCPServer(
    "Codebase Brain",
    "1.0.0",  // Version
    // ...
)
```

Update version during build:
```bash
VERSION=2.0.0 go build -o codebase-brain ./cmd/server
```

---

## Troubleshooting

### Go Not Found

Ensure Go is in your PATH:
```bash
export PATH=$PATH:/usr/local/go/bin
```

### Architecture Not Supported

Some older systems don't support certain architectures. Use:
- `GOARCH=amd64` for 64-bit x86
- `GOARCH=386` for 32-bit x86

### Cross-Compilation Errors

Some systems need CGO for certain features:
```bash
CGO_ENABLED=0 go build ...
```

---

## Output

Build produces a single executable:
- **Windows**: `codebase-brain.exe`
- **macOS/Linux**: `codebase-brain`

No additional dependencies required - static binary.

---

## Docker Build

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o codebase-brain ./cmd/server

FROM alpine:latest
COPY --from=builder /app/codebase-brain /usr/local/bin/
ENTRYPOINT ["codebase-brain"]
```

Build and run:
```bash
docker build -t codebase-brain .
docker run -v $(pwd):/project -e PROJECT_PATH=/project codebase-brain
```