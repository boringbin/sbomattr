// Package sbom provides internal utilities for working with SBOM (Software Bill of Materials) formats.
package sbom

import (
	"encoding/json"
	"errors"
)

// DetectFormat analyzes the SBOM data and returns the detected format string.
// It returns either "spdx" or "cyclonedx" based on format-specific markers in the JSON data.
// It supports both standard formats and GitHub-wrapped formats (e.g., {"sbom": {...}}).
func DetectFormat(data []byte) (string, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return "", err
	}

	// Check for GitHub wrapper format and unwrap if present
	if sbomData, hasWrapper := raw["sbom"]; hasWrapper {
		if sbomMap, ok := sbomData.(map[string]any); ok {
			raw = sbomMap
		}
	}

	// Check for SPDX markers
	if _, ok := raw["spdxVersion"]; ok {
		return "spdx", nil
	}
	if spdxID, ok := raw["SPDXID"].(string); ok && spdxID != "" {
		return "spdx", nil
	}

	// Check for CycloneDX markers
	if bomFormat, ok := raw["bomFormat"].(string); ok && bomFormat == "CycloneDX" {
		return "cyclonedx", nil
	}

	return "", errors.New("unknown SBOM format: could not detect SPDX or CycloneDX markers")
}
