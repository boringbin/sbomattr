package spdxextract

import (
	"encoding/json"
	"fmt"
)

// unwrapGitHubSBOM checks if the data is wrapped in GitHub's {"sbom": {...}} format and returns the unwrapped SPDX
// data if so, or the original data otherwise.
func unwrapGitHubSBOM(data []byte) ([]byte, error) {
	// Try to unmarshal as a map to check for GitHub wrapper
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Check for GitHub wrapper format: {"sbom": {...}}
	if sbomData, hasWrapper := wrapper["sbom"]; hasWrapper {
		return sbomData, nil
	}

	// Not wrapped, return original data
	return data, nil
}

// ParseSBOM parses SPDX JSON data from the given byte slice.
// It supports both standard SPDX format and GitHub-wrapped format ({"sbom": {...}}).
// It returns the parsed SPDX document or an error if parsing fails.
func ParseSBOM(data []byte) (*Document, error) {
	// Unwrap GitHub format if present
	unwrapped, err := unwrapGitHubSBOM(data)
	if err != nil {
		return nil, err
	}

	// Parse the JSON into an SPDX document
	var doc Document
	if unmarshalErr := json.Unmarshal(unwrapped, &doc); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to parse SBOM JSON: %w", unmarshalErr)
	}

	return &doc, nil
}
