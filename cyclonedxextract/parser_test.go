package cyclonedxextract_test

import (
	"testing"

	"github.com/boringbin/sbomattr/cyclonedxextract"
)

// TestParseSBOM_ValidJSON tests the ParseSBOM function with a valid JSON object.
func TestParseSBOM_ValidJSON(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{
		"bomFormat": "CycloneDX",
		"specVersion": "1.4",
		"components": [
			{
				"name": "lodash",
				"version": "4.17.21",
				"purl": "pkg:npm/lodash@4.17.21"
			}
		]
	}`)

	bom, err := cyclonedxextract.ParseSBOM(jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bom == nil {
		t.Fatal("Expected BOM, got nil")
	}

	if bom.BOMFormat != "CycloneDX" {
		t.Errorf("Expected BOMFormat 'CycloneDX', got %q", bom.BOMFormat)
	}

	if bom.SpecVersion != "1.4" {
		t.Errorf("Expected SpecVersion '1.4', got %q", bom.SpecVersion)
	}

	if len(bom.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(bom.Components))
	}

	if bom.Components[0].Name != "lodash" {
		t.Errorf("Expected component name 'lodash', got %q", bom.Components[0].Name)
	}
}

// TestParseSBOM_InvalidJSON tests the ParseSBOM function with an invalid JSON object.
func TestParseSBOM_InvalidJSON(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{this is not valid json}`)

	bom, err := cyclonedxextract.ParseSBOM(jsonData)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	if bom != nil {
		t.Errorf("Expected nil BOM for invalid JSON, got %+v", bom)
	}
}

// TestParseSBOM_EmptyJSON tests the ParseSBOM function with an empty JSON object.
func TestParseSBOM_EmptyJSON(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{}`)

	bom, err := cyclonedxextract.ParseSBOM(jsonData)
	if err != nil {
		t.Fatalf("Expected no error for empty JSON object, got %v", err)
	}

	if bom == nil {
		t.Fatal("Expected BOM, got nil")
	}

	if bom.BOMFormat != "" {
		t.Errorf("Expected empty BOMFormat, got %q", bom.BOMFormat)
	}

	if len(bom.Components) != 0 {
		t.Errorf("Expected 0 components, got %d", len(bom.Components))
	}
}

// TestParseSBOM_EmptyArray tests the ParseSBOM function with an empty components array.
func TestParseSBOM_EmptyArray(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{
		"bomFormat": "CycloneDX",
		"specVersion": "1.4",
		"components": []
	}`)

	bom, err := cyclonedxextract.ParseSBOM(jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bom == nil {
		t.Fatal("Expected BOM, got nil")
	}

	if len(bom.Components) != 0 {
		t.Errorf("Expected 0 components, got %d", len(bom.Components))
	}
}

// TestParseSBOM_ComplexLicense tests the ParseSBOM function with a complex license.
func TestParseSBOM_ComplexLicense(t *testing.T) {
	t.Parallel()

	jsonData := []byte(`{
		"bomFormat": "CycloneDX",
		"specVersion": "1.4",
		"components": [
			{
				"name": "test-package",
				"version": "1.0.0",
				"purl": "pkg:npm/test-package@1.0.0",
				"licenses": [
					{
						"license": {
							"id": "MIT",
							"name": "MIT License",
							"expression": "MIT OR Apache-2.0"
						}
					}
				]
			}
		]
	}`)

	bom, err := cyclonedxextract.ParseSBOM(jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bom == nil {
		t.Fatal("Expected BOM, got nil")
	}

	if len(bom.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(bom.Components))
	}

	component := bom.Components[0]
	if component.Licenses == nil {
		t.Fatal("Expected licenses to be set, got nil")
	}

	licenses := *component.Licenses
	if len(licenses) != 1 {
		t.Fatalf("Expected 1 license choice, got %d", len(licenses))
	}

	if licenses[0].License == nil {
		t.Fatal("Expected license to be set, got nil")
	}

	license := licenses[0].License
	if license.ID != "MIT" {
		t.Errorf("Expected license ID 'MIT', got %q", license.ID)
	}

	if license.Name != "MIT License" {
		t.Errorf("Expected license name 'MIT License', got %q", license.Name)
	}

	if license.Expression != "MIT OR Apache-2.0" {
		t.Errorf("Expected license expression 'MIT OR Apache-2.0', got %q", license.Expression)
	}
}

// TestParseSBOM_NullBytes tests the ParseSBOM function with null bytes.
func TestParseSBOM_NullBytes(t *testing.T) {
	t.Parallel()

	jsonData := []byte{}

	bom, err := cyclonedxextract.ParseSBOM(jsonData)
	if err == nil {
		t.Fatal("Expected error for empty bytes, got nil")
	}

	if bom != nil {
		t.Errorf("Expected nil BOM for empty bytes, got %+v", bom)
	}
}
