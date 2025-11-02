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
	"log/slog"
	"os"

	"github.com/boringbin/sbomattr/attribution"
	"github.com/boringbin/sbomattr/cyclonedxextract"
	"github.com/boringbin/sbomattr/internal/sbom"
	"github.com/boringbin/sbomattr/spdxextract"
)

// Process processes a single SBOM file provided as a byte slice.
// It automatically detects the SBOM format (SPDX or CycloneDX), parses it,
// and extracts attribution information.
//
// The context parameter can be used for cancellation.
// The logger parameter is optional; pass nil to disable logging.
//
// Returns a slice of Attribution structs or an error if the SBOM cannot be processed.
func Process(ctx context.Context, data []byte, logger *slog.Logger) ([]attribution.Attribution, error) {
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

	if logger != nil {
		logger.DebugContext(ctx, "detected SBOM format", "format", format)
	}

	// Extract attributions based on format
	switch format {
	case "spdx":
		doc, parseErr := spdxextract.ParseSBOM(data)
		if parseErr != nil {
			return nil, fmt.Errorf("parse SPDX: %w", parseErr)
		}
		return spdxextract.ExtractPackages(doc), nil
	case "cyclonedx":
		bom, parseErr := cyclonedxextract.ParseSBOM(data)
		if parseErr != nil {
			return nil, fmt.Errorf("parse CycloneDX: %w", parseErr)
		}
		return cyclonedxextract.ExtractPackages(bom), nil
	default:
		return nil, fmt.Errorf("unsupported SBOM format: %s", format)
	}
}

// ProcessFiles processes multiple SBOM files from the filesystem.
// It reads each file, processes the SBOM, aggregates the results, and deduplicates
// attributions based on Package URL (purl) or name if purl is not available.
//
// The context parameter can be used for cancellation.
// The logger parameter is optional; pass nil to disable logging.
// Errors processing individual files are logged but do not stop processing of other files.
//
// Returns the deduplicated attributions or an error if no valid attributions could be extracted.
func ProcessFiles(ctx context.Context, filenames []string, logger *slog.Logger) ([]attribution.Attribution, error) {
	var allAttributions []attribution.Attribution

	for _, filename := range filenames {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if logger != nil {
			logger.DebugContext(ctx, "processing file", "file", filename)
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			if logger != nil {
				logger.ErrorContext(ctx, "failed to read file", "file", filename, "error", err)
			}
			continue
		}

		attrs, err := Process(ctx, data, logger)
		if err != nil {
			if logger != nil {
				logger.ErrorContext(ctx, "failed to process file", "file", filename, "error", err)
			}
			continue
		}

		allAttributions = append(allAttributions, attrs...)
	}

	if len(allAttributions) == 0 {
		return nil, errors.New("no attributions extracted from any file")
	}

	// Deduplicate attributions
	deduplicated := attribution.Deduplicate(allAttributions, logger)

	return deduplicated, nil
}
