// Package sbomattr provides a high-level API for extracting attribution information
// from Software Bill of Materials (SBOM) files in SPDX and CycloneDX formats.
//
// Supported formats:
//   - SPDX 2.3 (JSON)
//   - CycloneDX 1.4 (JSON)
//   - GitHub-wrapped SBOMs (JSON)
package sbomattr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/boringbin/sbomattr/attribution"
	"github.com/boringbin/sbomattr/cyclonedxextract"
	"github.com/boringbin/sbomattr/internal/sbom"
	"github.com/boringbin/sbomattr/spdxextract"
)

// logger is the package-level logger. By default, logging is disabled.
//
//nolint:gochecknoglobals // Package-level logger is simpler
var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

// SetLogger configures the logger for this package.
// By default, logging is disabled. Call this function to enable logging.
func SetLogger(l *slog.Logger) {
	if l != nil {
		logger = l
	}
}

// Extractor defines the interface for SBOM format extractors.
type Extractor interface {
	Extract(data []byte) ([]attribution.Attribution, error)
}

// spdxExtractor implements Extractor for SPDX format.
type spdxExtractor struct{}

func (e *spdxExtractor) Extract(data []byte) ([]attribution.Attribution, error) {
	doc, err := spdxextract.ParseSBOM(data)
	if err != nil {
		return nil, fmt.Errorf("parse SPDX: %w", err)
	}
	return spdxextract.ExtractPackages(doc), nil
}

// cyclonedxExtractor implements Extractor for CycloneDX format.
type cyclonedxExtractor struct{}

func (e *cyclonedxExtractor) Extract(data []byte) ([]attribution.Attribution, error) {
	bom, err := cyclonedxextract.ParseSBOM(data)
	if err != nil {
		return nil, fmt.Errorf("parse CycloneDX: %w", err)
	}
	return cyclonedxextract.ExtractPackages(bom), nil
}

// Process processes a single SBOM file provided as a byte slice.
// It automatically detects the SBOM format (SPDX or CycloneDX), parses it,
// and extracts attribution information.
//
// The context parameter can be used for cancellation.
//
// Returns a slice of Attribution structs or an error if the SBOM cannot be processed.
func Process(ctx context.Context, data []byte) ([]attribution.Attribution, error) {
	// Check for cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Detect format
	format, err := sbom.DetectFormat(data)
	if err != nil {
		return nil, fmt.Errorf("detect format: %w", err)
	}

	logger.DebugContext(ctx, "detected SBOM format", "format", format) //nolint:sloglint // Package-level logger

	// Get the appropriate extractor
	extractors := map[string]Extractor{
		"spdx":      &spdxExtractor{},
		"cyclonedx": &cyclonedxExtractor{},
	}
	extractor, ok := extractors[format]
	if !ok {
		return nil, fmt.Errorf("unsupported SBOM format: %s", format)
	}

	// Extract attributions
	return extractor.Extract(data)
}

// ProcessFiles processes multiple SBOM files from the filesystem.
// It reads each file, processes the SBOM, aggregates the results, and deduplicates
// attributions based on Package URL (purl) or name if purl is not available.
//
// The context parameter can be used for cancellation.
// Errors processing individual files are logged but do not stop processing of other files.
//
// Returns the deduplicated attributions or an error if no valid attributions could be extracted.
func ProcessFiles(ctx context.Context, filenames []string) ([]attribution.Attribution, error) {
	var allAttributions []attribution.Attribution

	for _, filename := range filenames {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		logger.DebugContext(ctx, "processing file", "file", filename) //nolint:sloglint // Package-level logger

		data, err := os.ReadFile(filename)
		if err != nil {
			//nolint:sloglint // Package-level logger
			logger.ErrorContext(
				ctx,
				"failed to read file",
				"file",
				filename,
				"error",
				err,
			)
			continue
		}

		attrs, err := Process(ctx, data)
		if err != nil {
			//nolint:sloglint // Package-level logger
			logger.ErrorContext(
				ctx,
				"failed to process file",
				"file",
				filename,
				"error",
				err,
			)
			continue
		}

		allAttributions = append(allAttributions, attrs...)
	}

	if len(allAttributions) == 0 {
		return nil, errors.New("no attributions extracted from any file")
	}

	// Deduplicate attributions
	deduplicated := attribution.Deduplicate(allAttributions)

	return deduplicated, nil
}
