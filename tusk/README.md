# Tusk CLI

A command-line tool built with Go and Cobra.

## Installation

### From Source

```bash
git clone https://github.com/sboy99/projects.go/tusk.git
cd tusk
make build
```

### Using Go Install

```bash
go install github.com/sboy99/projects.go/tusk/cmd/tusk@latest
```

## Usage

```bash
# Show version information
tusk version

# Start a service
tusk start --service my-service

# Stop a service
tusk stop my-service

# Use a custom config file
tusk --config /path/to/config.yaml start --service my-service
```

## Configuration

Configuration can be provided via:
- Config file (YAML) - default location: `$HOME/.tusk/config.yaml`
- Environment variables (prefixed with `TUSK_`)
- Command-line flags

Example config file (`configs/config.yaml`):
```yaml
app_name: tusk
port: 8080
debug: false
```

## Development

### Build

```bash
make build
# or
./scripts/build.sh
```

### Run Tests

```bash
make test
# or
go test ./...
```

### Run E2E Tests

```bash
make test-e2e
# or
go test ./test/...
```

### Release

```bash
make release
# or
./scripts/release.sh
```

## Project Structure

```
tusk/
├── cmd/tusk/          # CLI entrypoint
├── internal/          # Private application code
│   ├── config/       # Configuration loader
│   ├── services/     # Core business logic
│   ├── executor/     # Command execution logic
│   ├── utils/        # Helper utilities
│   └── version/      # Version info
├── pkg/              # Public libraries
│   └── formatter/   # Formatting utilities
├── cli/              # Cobra command definitions
├── configs/          # Example configs
├── scripts/          # Build and release scripts
└── test/             # End-to-end tests
```

## License

MIT

