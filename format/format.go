package format

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/boringbin/sbomattr/attribution"
)

// CSV writes attributions as CSV to the provided io.Writer.
// The CSV has columns: Name, License, Purl, URL.
func CSV(w io.Writer, attributions []attribution.Attribution) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Name", "License", "Purl", "URL"}); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}

	// Write rows
	for _, a := range attributions {
		license := ""
		if a.License != nil {
			license = *a.License
		}

		url := ""
		if a.URL != nil {
			url = *a.URL
		}

		if err := writer.Write([]string{a.Name, license, a.Purl, url}); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}

	return nil
}

// JSON writes attributions as pretty-printed JSON to the provided io.Writer.
func JSON(w io.Writer, attributions []attribution.Attribution) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(attributions); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}
	return nil
}
