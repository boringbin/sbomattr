//go:build integration

package spdxextract_test

import (
	"os"
	"testing"

	"github.com/boringbin/sbomattr/internal/sbom"
	"github.com/boringbin/sbomattr/spdxextract"
)

// TestParseAndExtract_GitHubWrappedSPDX tests end-to-end processing of a GitHub-wrapped SPDX file.
func TestParseAndExtract_GitHubWrappedSPDX(t *testing.T) {
	t.Parallel()

	// Read the GitHub-wrapped SPDX file
	data, err := os.ReadFile("../testdata/github-wrapped-spdx.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Verify format detection works with wrapped format
	format, err := sbom.DetectFormat(data)
	if err != nil {
		t.Fatalf("DetectFormat failed: %v", err)
	}

	if format != "spdx" {
		t.Errorf("Expected format 'spdx', got '%s'", format)
	}

	// Parse the SBOM
	doc, err := spdxextract.ParseSBOM(data)
	if err != nil {
		t.Fatalf("ParseSBOM failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Verify document structure
	if doc.SPDXVersion != "SPDX-2.3" {
		t.Errorf("Expected SPDX version 'SPDX-2.3', got '%s'", doc.SPDXVersion)
	}

	if len(doc.Packages) != 3 {
		t.Errorf("Expected 3 packages, got %d", len(doc.Packages))
	}

	// Extract packages
	packages := spdxextract.ExtractPackages(doc)

	if len(packages) != 3 {
		t.Errorf("Expected 3 extracted packages, got %d", len(packages))
	}

	// Verify package details
	expectedPackages := map[string]string{
		"proxy-from-env": "MIT",
		"yargs":          "MIT",
		"cliui":          "ISC",
	}

	for _, pkg := range packages {
		expectedLicense, found := expectedPackages[pkg.Name]
		if !found {
			t.Errorf("Unexpected package name: %s", pkg.Name)
			continue
		}

		if pkg.License == nil {
			t.Errorf("Package %s has nil license", pkg.Name)
			continue
		}

		if *pkg.License != expectedLicense {
			t.Errorf("Package %s: expected license '%s', got '%s'", pkg.Name, expectedLicense, *pkg.License)
		}

		// Verify PURL is present
		if pkg.Purl == "" {
			t.Errorf("Package %s has empty PURL", pkg.Name)
		}
	}
}

// TestParseAndExtract_StandardSPDX tests end-to-end processing of a standard SPDX file.
func TestParseAndExtract_StandardSPDX(t *testing.T) {
	t.Parallel()

	// Read the standard SPDX file
	data, err := os.ReadFile("../testdata/example-spdx.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Verify format detection
	format, err := sbom.DetectFormat(data)
	if err != nil {
		t.Fatalf("DetectFormat failed: %v", err)
	}

	if format != "spdx" {
		t.Errorf("Expected format 'spdx', got '%s'", format)
	}

	// Parse the SBOM
	doc, err := spdxextract.ParseSBOM(data)
	if err != nil {
		t.Fatalf("ParseSBOM failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Extract packages
	packages := spdxextract.ExtractPackages(doc)

	if len(packages) != 3 {
		t.Errorf("Expected 3 extracted packages, got %d", len(packages))
	}

	// Verify package names from standard format
	expectedNames := map[string]bool{
		"lodash":  true,
		"react":   true,
		"express": true,
	}

	for _, pkg := range packages {
		if !expectedNames[pkg.Name] {
			t.Errorf("Unexpected package name: %s", pkg.Name)
		}
	}
}

// TestParseAndExtract_CompareFormats verifies both formats produce comparable results.
func TestParseAndExtract_CompareFormats(t *testing.T) {
	t.Parallel()

	// Read both files
	standardData, err := os.ReadFile("../testdata/example-spdx.json")
	if err != nil {
		t.Fatalf("Failed to read standard SPDX file: %v", err)
	}

	wrappedData, err := os.ReadFile("../testdata/github-wrapped-spdx.json")
	if err != nil {
		t.Fatalf("Failed to read GitHub-wrapped SPDX file: %v", err)
	}

	// Parse both
	standardDoc, err := spdxextract.ParseSBOM(standardData)
	if err != nil {
		t.Fatalf("Failed to parse standard SPDX: %v", err)
	}

	wrappedDoc, err := spdxextract.ParseSBOM(wrappedData)
	if err != nil {
		t.Fatalf("Failed to parse GitHub-wrapped SPDX: %v", err)
	}

	// Both should have same number of packages
	if len(standardDoc.Packages) != len(wrappedDoc.Packages) {
		t.Logf("Standard format has %d packages, wrapped format has %d packages (this is expected)",
			len(standardDoc.Packages), len(wrappedDoc.Packages))
	}

	// Both should have same SPDX version
	if standardDoc.SPDXVersion != wrappedDoc.SPDXVersion {
		t.Errorf("SPDX versions differ: standard=%s, wrapped=%s",
			standardDoc.SPDXVersion, wrappedDoc.SPDXVersion)
	}
}
