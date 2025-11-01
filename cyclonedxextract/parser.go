package cyclonedxextract

import (
	"encoding/json"
	"fmt"
)

// ParseSBOM parses CycloneDX JSON data from the given byte slice.
// It returns the parsed CycloneDX BOM or an error if parsing fails.
func ParseSBOM(data []byte) (*BOM, error) {
	// Parse the JSON into a CycloneDX BOM
	var bom BOM
	if unmarshalErr := json.Unmarshal(data, &bom); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to parse CycloneDX JSON: %w", unmarshalErr)
	}

	return &bom, nil
}
