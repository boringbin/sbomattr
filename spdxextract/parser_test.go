package spdxextract_test

import (
	"os"
	"testing"

	"github.com/boringbin/sbomattr/spdxextract"
)

// TestParseSBOM_StandardFormat tests parsing a standard SPDX JSON file.
func TestParseSBOM_StandardFormat(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../testdata/example-spdx.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	doc, err := spdxextract.ParseSBOM(data)
	if err != nil {
		t.Fatalf("ParseSBOM failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	if doc.SPDXVersion != "SPDX-2.3" {
		t.Errorf("Expected SPDX version 'SPDX-2.3', got '%s'", doc.SPDXVersion)
	}

	const expectedPackages = 3
	if len(doc.Packages) != expectedPackages {
		t.Errorf("Expected %d packages, got %d", expectedPackages, len(doc.Packages))
	}
}

// TestParseSBOM_GitHubWrappedFormat tests parsing a GitHub-wrapped SPDX JSON file.
func TestParseSBOM_GitHubWrappedFormat(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../testdata/github-wrapped-spdx.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	doc, err := spdxextract.ParseSBOM(data)
	if err != nil {
		t.Fatalf("ParseSBOM failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	if doc.SPDXVersion != "SPDX-2.3" {
		t.Errorf("Expected SPDX version 'SPDX-2.3', got '%s'", doc.SPDXVersion)
	}

	const expectedPackages = 3
	if len(doc.Packages) != expectedPackages {
		t.Errorf("Expected %d packages, got %d", expectedPackages, len(doc.Packages))
	}

	// Verify specific package from GitHub format
	if doc.Packages[0].Name != "proxy-from-env" {
		t.Errorf("Expected first package name 'proxy-from-env', got '%s'", doc.Packages[0].Name)
	}
}

// TestParseSBOM_InvalidJSON tests that invalid JSON returns an error.
func TestParseSBOM_InvalidJSON(t *testing.T) {
	t.Parallel()

	invalidData := []byte("not valid json")

	_, err := spdxextract.ParseSBOM(invalidData)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}

// TestParseSBOM_EmptyWrapper tests parsing JSON with an empty sbom wrapper.
func TestParseSBOM_EmptyWrapper(t *testing.T) {
	t.Parallel()

	data := []byte(`{"sbom": {}}`)

	doc, err := spdxextract.ParseSBOM(data)
	if err != nil {
		t.Fatalf("ParseSBOM failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected non-nil document")
	}

	// Empty document should have no packages
	if len(doc.Packages) != 0 {
		t.Errorf("Expected 0 packages, got %d", len(doc.Packages))
	}
}
