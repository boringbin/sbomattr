package cyclonedxextract_test

import (
	"testing"

	"github.com/boringbin/sbomattr/cyclonedxextract"
)

// TestExtractPackages_NilBOM tests the ExtractPackages function with a nil BOM.
func TestExtractPackages_NilBOM(t *testing.T) {
	t.Parallel()

	result := cyclonedxextract.ExtractPackages(nil)

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(result))
	}
}

// TestExtractPackages_EmptyComponents tests the ExtractPackages function with an empty components slice.
func TestExtractPackages_EmptyComponents(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components:  []cyclonedxextract.Component{},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(result))
	}
}

// TestExtractPackages_NilComponentsSlice tests the ExtractPackages function with a nil components slice.
func TestExtractPackages_NilComponentsSlice(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components:  nil,
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(result))
	}
}

// TestExtractPackages_ComponentWithAllFields tests the ExtractPackages function with a component with all fields.
func TestExtractPackages_ComponentWithAllFields(t *testing.T) {
	t.Parallel()

	mitLicense := "MIT"
	licenses := cyclonedxextract.Licenses{
		{
			License: &cyclonedxextract.License{
				ID: mitLicense,
			},
		},
	}

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "lodash",
				Version:  "4.17.21",
				Purl:     "pkg:npm/lodash@4.17.21",
				Licenses: &licenses,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Name != "lodash" {
		t.Errorf("Expected name 'lodash', got %q", attr.Name)
	}

	if attr.Purl != "pkg:npm/lodash@4.17.21" {
		t.Errorf("Expected purl 'pkg:npm/lodash@4.17.21', got %q", attr.Purl)
	}

	if attr.License == nil {
		t.Fatal("Expected license to be set, got nil")
	}

	if *attr.License != "MIT" {
		t.Errorf("Expected license 'MIT', got %q", *attr.License)
	}

	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	expectedURL := "https://www.npmjs.com/package/lodash/v/4.17.21"
	if *attr.URL != expectedURL {
		t.Errorf("Expected URL %q, got %q", expectedURL, *attr.URL)
	}
}

// TestExtractPackages_ComponentWithoutPurl tests the ExtractPackages function with a component without a purl.
func TestExtractPackages_ComponentWithoutPurl(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:    "custom-package",
				Version: "1.0.0",
				Purl:    "", // No purl
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Name != "custom-package" {
		t.Errorf("Expected name 'custom-package', got %q", attr.Name)
	}

	if attr.Purl != "" {
		t.Errorf("Expected empty purl, got %q", attr.Purl)
	}

	if attr.URL != nil {
		t.Errorf("Expected URL to be nil, got %q", *attr.URL)
	}
}

// TestExtractPackages_ComponentWithoutLicense tests the ExtractPackages function with a component without a license.
func TestExtractPackages_ComponentWithoutLicense(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "unlicensed-package",
				Version:  "1.0.0",
				Purl:     "pkg:npm/unlicensed-package@1.0.0",
				Licenses: nil,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.License != nil {
		t.Errorf("Expected license to be nil, got %q", *attr.License)
	}
}

// TestExtractPackages_MultipleComponents tests the ExtractPackages function with multiple components.
func TestExtractPackages_MultipleComponents(t *testing.T) {
	t.Parallel()

	apacheLicense := "Apache-2.0"
	mitLicense := "MIT"
	licenses1 := cyclonedxextract.Licenses{
		{
			License: &cyclonedxextract.License{
				ID: apacheLicense,
			},
		},
	}
	licenses2 := cyclonedxextract.Licenses{
		{
			License: &cyclonedxextract.License{
				ID: mitLicense,
			},
		},
	}

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "package-one",
				Version:  "1.0.0",
				Purl:     "pkg:npm/package-one@1.0.0",
				Licenses: &licenses1,
			},
			{
				Name:     "package-two",
				Version:  "2.0.0",
				Purl:     "pkg:pypi/package-two@2.0.0",
				Licenses: &licenses2,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 2 {
		t.Fatalf("Expected 2 attributions, got %d", len(result))
	}

	if result[0].Name != "package-one" {
		t.Errorf("Expected first package name 'package-one', got %q", result[0].Name)
	}

	if result[1].Name != "package-two" {
		t.Errorf("Expected second package name 'package-two', got %q", result[1].Name)
	}
}

// TestExtractLicense_NilLicenses tests the ExtractPackages function with a nil licenses slice.
func TestExtractLicense_NilLicenses(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "test-package",
				Version:  "1.0.0",
				Purl:     "pkg:npm/test-package@1.0.0",
				Licenses: nil,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	if result[0].License != nil {
		t.Errorf("Expected nil license, got %q", *result[0].License)
	}
}

// TestExtractLicense_EmptyLicensesArray tests the ExtractPackages function with an empty licenses array.
func TestExtractLicense_EmptyLicensesArray(t *testing.T) {
	t.Parallel()

	emptyLicenses := cyclonedxextract.Licenses{}
	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "test-package",
				Version:  "1.0.0",
				Purl:     "pkg:npm/test-package@1.0.0",
				Licenses: &emptyLicenses,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	if result[0].License != nil {
		t.Errorf("Expected nil license, got %q", *result[0].License)
	}
}

// TestExtractLicense_WithExpression tests the ExtractPackages function with a license expression.
func TestExtractLicense_WithExpression(t *testing.T) {
	t.Parallel()

	expression := "MIT OR Apache-2.0"
	licenses := cyclonedxextract.Licenses{
		{
			License: &cyclonedxextract.License{
				Expression: expression,
				ID:         "MIT",   // Should be ignored in favor of expression
				Name:       "Other", // Should be ignored in favor of expression
			},
		},
	}

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "test-package",
				Version:  "1.0.0",
				Purl:     "pkg:npm/test-package@1.0.0",
				Licenses: &licenses,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	if result[0].License == nil {
		t.Fatal("Expected license to be set, got nil")
	}

	if *result[0].License != expression {
		t.Errorf("Expected license %q, got %q", expression, *result[0].License)
	}
}

// TestExtractLicense_WithIDOnly tests the ExtractPackages function with a license ID only.
func TestExtractLicense_WithIDOnly(t *testing.T) {
	t.Parallel()

	licenseID := "Apache-2.0"
	licenses := cyclonedxextract.Licenses{
		{
			License: &cyclonedxextract.License{
				ID:   licenseID,
				Name: "Other", // Should be ignored in favor of ID
			},
		},
	}

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "test-package",
				Version:  "1.0.0",
				Purl:     "pkg:npm/test-package@1.0.0",
				Licenses: &licenses,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	if result[0].License == nil {
		t.Fatal("Expected license to be set, got nil")
	}

	if *result[0].License != licenseID {
		t.Errorf("Expected license %q, got %q", licenseID, *result[0].License)
	}
}

// TestExtractLicense_WithNameOnly tests the ExtractPackages function with a license name only.
func TestExtractLicense_WithNameOnly(t *testing.T) {
	t.Parallel()

	licenseName := "BSD-3-Clause"
	licenses := cyclonedxextract.Licenses{
		{
			License: &cyclonedxextract.License{
				Name: licenseName,
			},
		},
	}

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "test-package",
				Version:  "1.0.0",
				Purl:     "pkg:npm/test-package@1.0.0",
				Licenses: &licenses,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	if result[0].License == nil {
		t.Fatal("Expected license to be set, got nil")
	}

	if *result[0].License != licenseName {
		t.Errorf("Expected license %q, got %q", licenseName, *result[0].License)
	}
}

// TestExtractLicense_WithNilLicenseField tests the ExtractPackages function with a nil license field.
func TestExtractLicense_WithNilLicenseField(t *testing.T) {
	t.Parallel()

	licenses := cyclonedxextract.Licenses{
		{
			License: nil,
		},
	}

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:     "test-package",
				Version:  "1.0.0",
				Purl:     "pkg:npm/test-package@1.0.0",
				Licenses: &licenses,
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	if result[0].License != nil {
		t.Errorf("Expected nil license, got %q", *result[0].License)
	}
}
