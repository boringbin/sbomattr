package spdxextract

import (
	"github.com/boringbin/sbomattr/attribution"
)

// ExtractPackages extracts a simplified list of packages from an SPDX document.
// It returns a slice of Attribution structs containing name, version, purl, and license information.
func ExtractPackages(doc *Document) []attribution.Attribution {
	if doc == nil || doc.Packages == nil {
		return []attribution.Attribution{}
	}

	packages := make([]attribution.Attribution, 0, len(doc.Packages))

	for _, pkg := range doc.Packages {
		// Prefer concluded license, fall back to declared license
		license := pkg.LicenseConcluded
		if license == "" || license == "NOASSERTION" {
			license = pkg.LicenseDeclared
		}

		p := attribution.Attribution{
			Name:    pkg.Name,
			License: &license,
		}

		// Extract purl from external references
		for _, ref := range pkg.ExternalRefs {
			if ref.ReferenceType == "purl" {
				p.Purl = ref.ReferenceLocator
				break
			}
		}

		// Construct URL from purl if not already present
		if p.URL == nil && p.Purl != "" {
			p.URL = attribution.PurlToURL(p.Purl)
		}

		packages = append(packages, p)
	}

	return packages
}
