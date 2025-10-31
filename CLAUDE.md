# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

sbomattr is a CLI tool written in Go that creates aggregated notices for one or more SBOMs (Software Bill of Materials) in SPDX or CycloneDX formats.

**Key Features:**
- **Multi-format support**: SPDX 2.3 and CycloneDX 1.4 JSON formats
- **GitHub integration**: Handles GitHub-wrapped SBOM format (`{"sbom": {...}}`)
- **Automatic format detection**: Detects SPDX vs CycloneDX automatically
- **Aggregation**: Process multiple SBOM files in a single command
- **Deduplication**: Removes duplicate packages by purl or name
- **Package URL conversion**: Converts purls to package manager URLs
- **Multiple output formats**: CSV (default) and JSON
- **Context-aware**: Full context.Context support for cancellation
- **Structured logging**: Uses log/slog with configurable verbosity

**Supported Package URL (purl) Types:**
cargo, composer, gem, golang, maven, npm, nuget, pub, pypi, github (with URL conversion support)

## CLI Usage

```bash
# Basic usage - single file
./bin/sbomattr sbom.json

# Multiple files - aggregates and deduplicates
./bin/sbomattr sbom1.json sbom2.json sbom3.json

# Directory expansion - processes all .json files
./bin/sbomattr ./sboms/

# Verbose logging
./bin/sbomattr -v sbom.json

# Check version
./bin/sbomattr -version
```

**Output:** CSV format to stdout (Name, License, Purl, URL columns)

**Exit Codes:**
- 0: Success
- 1: Invalid arguments
- 2: Invalid SBOM format
- 3: Runtime error

## Development Commands

### Build
```bash
make all           # Build the binary to bin/sbomattr
go build -o bin/sbomattr .
```

### Testing
```bash
make test          # Run unit tests (excludes integration tests, uses -short flag)
make test-integration  # Run only integration tests (requires -tags=integration)
make test-all      # Run all tests including integration
make test-coverage # Run tests with coverage report
```

### Code Quality
```bash
make check         # Run format-check and lint-check
make fix           # Run format-fix and lint-fix

make format-check  # Verify code is formatted with gofmt
make format-fix    # Format code with gofmt

make lint-check    # Run golangci-lint
make lint-fix      # Run golangci-lint with --fix

make vet           # Run go vet
```

### Maintenance
```bash
make tidy          # Run go mod tidy
make clean         # Remove bin/sbomattr, coverage.out, coverage.html
```

## Code Standards

### Linting Configuration

The project uses an extremely strict golangci-lint configuration (.golangci.yaml) based on [maratori's golden config](https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322). Key requirements:

- **Line length**: Maximum 120 characters (enforced by golines formatter)
- **Import formatting**: Uses goimports with local prefix `github.com/boringbin/sbomattr`
- **Cyclomatic complexity**: Max 30 per function, average 10 per package
- **Function size**: Max 100 lines, max 50 statements
- **Cognitive complexity**: Max 20
- **No naked returns**: Naked returns are disallowed (max-func-lines: 0)
- **Logging**: Must use log/slog (not log package) except in main.go
- **Random numbers**: Must use math/rand/v2 (not math/rand) in non-test files

### Strict Rules Enforced

- No global variables (gochecknoglobals)
- No init functions (gochecknoinits)
- All errors must be checked (errcheck with type assertions)
- Exhaustive switch/map checks for enums
- No magic numbers (mnd linter)
- Context must be used for HTTP requests (noctx)
- Security checks enabled (gosec)
- Proper struct field tags required (musttag)

### Testing Standards

- Unit tests: Use `-short` flag to skip integration tests
- Integration tests: Require `-tags=integration` build tag
- Test files: Relaxed linting (no bodyclose, funlen, gosec, etc. in _test.go)
- Separate test packages: Use `testpackage` linter to encourage `_test` package suffix
- Parallel testing: Use `t.Parallel()` appropriately (tparallel linter)

## Project Structure

```
sbomattr/
├── .devcontainer/              # VS Code devcontainer (Go 1.25, golangci-lint 2.5.0)
├── .github/
│   ├── renovate.json           # Renovate bot configuration
│   └── workflows/              # CI/CD workflows
│       ├── build.yaml          # Build validation
│       ├── test.yaml           # Test coverage
│       └── typos.yaml          # Spell checking
├── attribution/                # Core data types and utilities
│   ├── attribution.go          # Attribution struct definition
│   ├── deduplicate.go          # Deduplication logic (by purl/name)
│   ├── url.go                  # Purl to package manager URL conversion
│   └── *_test.go               # Tests
├── cmd/sbomattr/               # CLI entry point
│   └── main.go                 # Main CLI implementation
├── cyclonedxextract/           # CycloneDX SBOM parser
│   ├── types.go                # CycloneDX type definitions
│   ├── parser.go               # JSON parsing
│   └── extractor.go            # Attribution extraction
├── format/                     # Output formatters
│   ├── format.go               # CSV and JSON formatters
│   └── format_test.go          # Tests
├── internal/sbom/              # Internal SBOM utilities
│   ├── sbom.go                 # Format detection (SPDX/CycloneDX/GitHub)
│   └── sbom_test.go            # Tests
├── spdxextract/                # SPDX SBOM parser
│   ├── types.go                # SPDX type definitions
│   ├── parser.go               # JSON parsing with GitHub wrapper support
│   ├── parser_integration_test.go  # Integration tests (-tags=integration)
│   └── extractor.go            # Attribution extraction
├── testdata/                   # Test fixtures
│   ├── example-spdx.json       # SPDX 2.3 example (npm packages)
│   ├── example-cyclonedx.json  # CycloneDX 1.4 example (pypi/npm)
│   └── github-wrapped-spdx.json # GitHub-wrapped SPDX format
├── bin/                        # Build output directory
├── processor.go                # Main processing logic (root package)
├── processor_test.go           # Processor tests
├── .golangci.yaml              # Extremely strict linting configuration
├── Makefile                    # Build automation
├── go.mod                      # Go 1.25.0
└── typos.toml                  # Typos checker configuration
```

## Package Architecture

### Root Package (`github.com/boringbin/sbomattr`)
**Files:** `processor.go`, `processor_test.go`

**Purpose:** Main processing logic and extractor orchestration

**Key Types:**
```go
type Extractor interface {
    Extract(ctx context.Context, data []byte) ([]attribution.Attribution, error)
}
```

**Key Functions:**
- `Process(ctx context.Context, data []byte) ([]attribution.Attribution, error)` - Process single SBOM
- `ProcessFiles(ctx context.Context, filenames []string) ([]attribution.Attribution, error)` - Process multiple files with aggregation/deduplication
- `SetLogger(*slog.Logger)` - Configure package-level logger

**Features:**
- Automatic format detection via `internal/sbom.DetectFormat`
- Extractor pattern for format-agnostic processing
- Context support for cancellation
- Package-level logger (disabled by default)

### `attribution` Package
**Purpose:** Core data types and utilities

**Key Type:**
```go
type Attribution struct {
    Name    string  // Package name
    License *string // Declared license (optional, pointer for nil distinction)
    URL     *string // Package URL (optional, pointer for nil distinction)
    Purl    string  // Package URL (purl format)
}
```

**Key Functions:**
- `Deduplicate([]Attribution) []Attribution` - Remove duplicates by purl (fallback to name)
- `PurlToURL(purlString string) (*string, error)` - Convert purl to package manager URL
- `SetLogger(*slog.Logger)` - Configure package-level logger

**Purl URL Conversion Support:**
- cargo → crates.io
- composer → packagist.org
- gem → rubygems.org
- golang → pkg.go.dev
- maven → mvnrepository.com
- npm → npmjs.com (handles org packages)
- nuget → nuget.org
- pub → pub.dev
- pypi → pypi.org
- github → github.com

### `cmd/sbomattr` Package
**Purpose:** CLI entry point and argument parsing

**Key Features:**
- Flag parsing: `-v` (verbose), `-version`
- Path expansion (files and directories)
- JSON file filtering (`.json` extension)
- Exit code handling (0=success, 1=invalid args, 2=invalid SBOM, 3=runtime error)
- Logger configuration based on verbosity
- CSV output to stdout

### `cyclonedxextract` Package
**Purpose:** Parse and extract attribution from CycloneDX SBOMs

**Key Types:**
- `BOM` - Top-level CycloneDX structure
- `Component` - Component with name, version, purl, licenses
- `License`, `LicenseChoice`, `Licenses` - License representation

**Key Functions:**
- `ParseSBOM([]byte) (*BOM, error)` - Parse CycloneDX JSON
- `ExtractPackages(*BOM) []attribution.Attribution` - Extract attributions

**License Extraction Priority:** expression > ID > name

### `format` Package
**Purpose:** Output formatters for attribution data

**Key Functions:**
- `CSV(w io.Writer, attrs []attribution.Attribution) error` - CSV output with header
- `JSON(w io.Writer, attrs []attribution.Attribution) error` - Pretty-printed JSON (2-space indent)

**CSV Format:** Name, License, Purl, URL (with proper escaping)

### `internal/sbom` Package
**Purpose:** Internal SBOM utilities (not importable externally)

**Key Function:**
- `DetectFormat(data []byte) (string, error)` - Auto-detect SBOM format

**Supported Formats:**
- SPDX (checks for `spdxVersion`, `SPDXID`)
- CycloneDX (checks for `bomFormat: "CycloneDX"`)
- GitHub-wrapped (checks for `{"sbom": {...}}` wrapper)

### `spdxextract` Package
**Purpose:** Parse and extract attribution from SPDX SBOMs

**Key Types:**
- `Document` - SPDX document with version, SPDXID, packages
- `Package` - Package with name, version, licenses, external refs
- `ExternalRef` - External references (including purl)

**Key Functions:**
- `ParseSBOM([]byte) (*Document, error)` - Parse SPDX JSON with GitHub wrapper support
- `unwrapGitHubSBOM([]byte) ([]byte, error)` - Handle GitHub's `{"sbom": {...}}` format
- `ExtractPackages(*Document) []attribution.Attribution` - Extract attributions

**License Extraction:** Prefers concluded license, falls back to declared if "NOASSERTION"

## Dependencies

**External:**
- `github.com/package-url/packageurl-go v0.1.3` - Purl parsing and validation

**Standard Library:**
- `context` - Cancellation support
- `encoding/csv`, `encoding/json` - Parsing and formatting
- `log/slog` - Structured logging
- `flag` - CLI argument parsing
- `os`, `io`, `path/filepath` - File operations

## Test Data

The `testdata/` directory contains example SBOMs for testing:

1. **example-spdx.json** - SPDX 2.3 format
   - 3 npm packages: lodash@4.17.21, react@18.2.0, express@4.18.2
   - Demonstrates concluded vs declared licenses
   - Express has "NOASSERTION" for concluded license

2. **example-cyclonedx.json** - CycloneDX 1.4 format
   - 4 packages: requests, numpy, flask (pypi), lodash (npm)
   - Various license formats: ID, name, expression
   - Demonstrates different license structures

3. **github-wrapped-spdx.json** - GitHub-wrapped SPDX format
   - GitHub's `{"sbom": {...}}` wrapper
   - 3 packages: proxy-from-env, yargs, cliui
   - Generated by GitHub Dependency Graph and protobom

## Key Design Patterns

1. **Package-level Loggers:** Both root package and attribution package use package-level slog loggers (disabled by default, configurable via SetLogger). This avoids global variables while maintaining linter compliance.

2. **Interface-based Extraction:** `Extractor` interface allows format-agnostic processing in the main processor.

3. **Pointer Fields for Optional Data:** `License` and `URL` fields use pointers (`*string`) to distinguish between empty string and missing data.

4. **Deduplication Strategy:** Primary key is purl, falls back to name if purl is empty.

5. **Context Propagation:** All processing functions accept `context.Context` for cancellation support.

6. **Test Separation:** Integration tests use build tags (`//go:build integration`) and are excluded from normal test runs.

7. **Parallel Testing:** Tests use `t.Parallel()` extensively to improve test execution speed.

8. **Error Wrapping:** Consistent use of `fmt.Errorf` with `%w` for error context.

## GitHub Actions CI/CD

**Workflows:**
- **build.yaml** - Build validation (runs `make` on push/PR to main)
- **test.yaml** - Test coverage (runs `make test-coverage` on push/PR to main)
- **typos.yaml** - Spell checking with crate-ci/typos v1.38.1

**Environment:**
- Go 1.25.0
- Uses commit hash for checkout action (security best practice)

## Development Environment

**.devcontainer:**
- Base: VS Code Go devcontainer with Go 1.25
- Includes: golangci-lint v2.5.0
- Features: Debugging support (SYS_PTRACE capability)

**typos.toml:**
- Checks: CLAUDE.md, LICENSE, README.md, *.go files
- Excludes: All other files

## Go Version

Uses Go 1.25.0 (specified in go.mod)
