# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## Design Philosophy

This project follows a **minimal and simple** design philosophy. The smaller the public API surface, the easier it is to maintain.

## Project Overview

sbomattr is a CLI tool written in Go that creates aggregated notices for SBOMs (Software Bill of Materials) in SPDX or CycloneDX formats.

**Features:**
- Multi-format support (SPDX 2.3, CycloneDX 1.4 JSON)
- GitHub-wrapped SBOM support (`{"sbom": {...}}`)
- Automatic format detection
- Aggregation and deduplication by purl/name
- Package URL conversion (29 purl types: cargo, composer, gem, golang, maven, npm, nuget, pub, pypi, github, etc.)
- CSV (default) and JSON output
- Context-aware with structured logging

## CLI Usage

```bash
./bin/sbomattr sbom.json                      # Single file
./bin/sbomattr sbom1.json sbom2.json          # Multiple files (aggregates)
./bin/sbomattr ./sboms/                       # Directory (all .json files)
./bin/sbomattr -v sbom.json                   # Verbose logging
./bin/sbomattr -version                       # Check version
```

**Output:** CSV to stdout (Name, License, Purl, URL)

**Exit Codes:**
- 0: Success
- 1: Invalid arguments
- 2: Invalid SBOM format
- 3: Runtime error

## Development Commands

```bash
# Build
make all                  # Build to bin/sbomattr
go build -o bin/sbomattr .

# Test
make test                 # Unit tests (-short flag)
make test-integration     # Integration tests (-tags=integration)
make test-all             # All tests
make test-coverage        # Coverage report

# Quality
make check                # Run format-check and lint-check
make fix                  # Run format-fix and lint-fix
make vet                  # Run go vet

# Maintenance
make tidy                 # Run go mod tidy
make clean                # Remove bin/, coverage files
```

## Architecture

### Package Structure

```
sbomattr/
├── attribution/          # Core types, deduplication, purl→URL conversion
├── cmd/sbomattr/         # CLI entry point
├── cyclonedxextract/     # CycloneDX parser
├── spdxextract/          # SPDX parser (with GitHub wrapper support)
├── format/               # CSV and JSON formatters
├── internal/sbom/        # Format detection
├── testdata/             # Test fixtures
└── processor.go          # Root package: main processing logic
```

### Key Types and Functions

**Root package** (`github.com/boringbin/sbomattr`):
```go
Process(ctx context.Context, data []byte, logger *slog.Logger) ([]Attribution, error)
ProcessFiles(ctx context.Context, filenames []string, logger *slog.Logger) ([]Attribution, error)
```

**attribution package**:
```go
type Attribution struct {
    Name    string   // Package name
    License *string  // Optional (pointer for nil vs empty)
    URL     *string  // Optional (pointer for nil vs empty)
    Purl    string   // Package URL
}

Deduplicate(attributions []Attribution, logger *slog.Logger) []Attribution
PurlToURL(purlString string, logger *slog.Logger) (*string, error)
```

**Sentinel errors**:
- `attribution.ErrEmptyPurl` - Empty/whitespace purl string
- `attribution.ErrUnsupportedPurlType` - Unsupported purl type

**URL preference**: SBOM-provided URL > purl-generated URL

**Format packages**:
- `cyclonedxextract.ParseSBOM(data) (*BOM, error)` + `ExtractPackages(bom)`
- `spdxextract.ParseSBOM(data) (*Document, error)` + `ExtractPackages(doc)`
- `format.CSV(w, attrs)` and `format.JSON(w, attrs)`

## Code Standards

- **Linting**: Strict golangci-lint config (.golangci.yaml)
- **Line length**: Max 120 characters
- **Logging**: Use log/slog (not log package)
- **Testing**: Use `t.Parallel()` except for tests modifying global state (CLI tests)
- **Integration tests**: Require `-tags=integration` build tag
- **Error handling**: Use sentinel errors for expected conditions, wrap errors with `%w`

## Dependencies

**External**: `github.com/package-url/packageurl-go v0.1.3`

**Standard library**: context, encoding/csv, encoding/json, log/slog, flag, os, io, path/filepath

## Key Patterns

1. **Explicit logger parameters**: Pass `*slog.Logger` to functions (nil to disable logging)
2. **Direct dispatch**: Simple `switch` statement for format selection (no interface abstraction)
3. **Pointer fields**: `*string` for optional data (nil vs empty distinction)
4. **Context propagation**: All processing functions accept `context.Context`
5. **Deduplication**: Primary key is purl, fallback to name
6. **Modern Go**: Uses `any` instead of `interface{}`, explicit error returns

## Environment

- Go 1.25.0
- GitHub Actions: build, test (with Codecov), CodeQL security analysis, typos checker
- Devcontainer: Go 1.25 + golangci-lint v2.5.0
