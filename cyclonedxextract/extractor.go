package cyclonedxextract

import (
	"github.com/boringbin/sbomattr/attribution"
)

// ExtractPackages extracts a simplified list of packages from a CycloneDX BOM.
// It returns a slice of Attribution structs containing name, version, purl, and license information.
func ExtractPackages(bom *BOM) []attribution.Attribution {
	if bom == nil || bom.Components == nil {
		return []attribution.Attribution{}
	}

	packages := make([]attribution.Attribution, 0, len(bom.Components))

	for _, component := range bom.Components {
		p := attribution.Attribution{
			Name: component.Name,
		}

		// Extract purl if available
		if component.Purl != "" {
			p.Purl = component.Purl
		}

		// Construct URL: prefer external references, fall back to purl conversion
		if refURL := findBestExternalRefURL(component.ExternalReferences); refURL != nil {
			p.URL = refURL
		} else if p.Purl != "" {
			p.URL = attribution.PurlToURL(p.Purl)
		}

		// Extract license information
		if component.Licenses != nil {
			license := extractLicense(component.Licenses)
			if license != nil {
				p.License = license
			}
		}

		packages = append(packages, p)
	}

	return packages
}

// extractLicense extracts license information from CycloneDX Licenses structure.
// It prefers license expressions, then license IDs, then license names.
func extractLicense(licenses *Licenses) *string {
	if licenses == nil || len(*licenses) == 0 {
		return nil
	}

	firstLicense := (*licenses)[0]

	// Prefer expression (e.g., "MIT OR Apache-2.0")
	if firstLicense.License != nil {
		if firstLicense.License.Expression != "" {
			return &firstLicense.License.Expression
		}

		// Fall back to License ID or Name
		if firstLicense.License.ID != "" {
			return &firstLicense.License.ID
		}
		if firstLicense.License.Name != "" {
			return &firstLicense.License.Name
		}
	}

	return nil
}

// findBestExternalRefURL finds the best URL from external references.
// Priority order: website > distribution > documentation > vcs.
func findBestExternalRefURL(refs []ExternalReference) *string {
	if len(refs) == 0 {
		return nil
	}

	// Priority order for reference types
	priorityOrder := []string{"website", "distribution", "documentation", "vcs"}

	for _, refType := range priorityOrder {
		for _, ref := range refs {
			if ref.Type == refType && ref.URL != "" {
				return &ref.URL
			}
		}
	}

	return nil
}
