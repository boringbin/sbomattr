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

// TestExtractPackages_WithExternalRefWebsite tests that "website" external ref is preferred over purl.
func TestExtractPackages_WithExternalRefWebsite(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:    "requests",
				Version: "2.28.1",
				Purl:    "pkg:pypi/requests@2.28.1",
				ExternalReferences: []cyclonedxextract.ExternalReference{
					{
						Type: "website",
						URL:  "https://requests.readthedocs.io/",
					},
				},
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Website external ref should be preferred over purl-generated URL
	if *attr.URL != "https://requests.readthedocs.io/" {
		t.Errorf("Expected URL to be website ref 'https://requests.readthedocs.io/', got %q", *attr.URL)
	}
}

// TestExtractPackages_WithMultipleExternalRefs tests priority order of external refs.
func TestExtractPackages_WithMultipleExternalRefs(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:    "flask",
				Version: "2.3.0",
				Purl:    "pkg:pypi/flask@2.3.0",
				ExternalReferences: []cyclonedxextract.ExternalReference{
					{
						Type: "vcs",
						URL:  "https://github.com/pallets/flask",
					},
					{
						Type: "documentation",
						URL:  "https://flask.palletsprojects.com/",
					},
					{
						Type: "website",
						URL:  "https://palletsprojects.com/p/flask/",
					},
				},
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Website should be preferred (highest priority)
	if *attr.URL != "https://palletsprojects.com/p/flask/" {
		t.Errorf("Expected URL to be website ref (highest priority), got %q", *attr.URL)
	}
}

// TestExtractPackages_WithExternalRefVCS tests that "vcs" external ref is used when website is not available.
func TestExtractPackages_WithExternalRefVCS(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:    "numpy",
				Version: "1.24.0",
				Purl:    "pkg:pypi/numpy@1.24.0",
				ExternalReferences: []cyclonedxextract.ExternalReference{
					{
						Type: "vcs",
						URL:  "https://github.com/numpy/numpy",
					},
				},
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// VCS external ref should be used
	if *attr.URL != "https://github.com/numpy/numpy" {
		t.Errorf("Expected URL to be vcs ref 'https://github.com/numpy/numpy', got %q", *attr.URL)
	}
}

// TestExtractPackages_WithExternalRefDistribution tests that "distribution" ref is preferred over "documentation".
func TestExtractPackages_WithExternalRefDistribution(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:    "django",
				Version: "4.2.0",
				Purl:    "pkg:pypi/django@4.2.0",
				ExternalReferences: []cyclonedxextract.ExternalReference{
					{
						Type: "documentation",
						URL:  "https://docs.djangoproject.com/",
					},
					{
						Type: "distribution",
						URL:  "https://www.djangoproject.com/download/",
					},
				},
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Distribution should be preferred over documentation
	if *attr.URL != "https://www.djangoproject.com/download/" {
		t.Errorf("Expected URL to be distribution ref, got %q", *attr.URL)
	}
}

// TestExtractPackages_WithExternalRefEmptyURL tests that empty URL falls back to purl.
func TestExtractPackages_WithExternalRefEmptyURL(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:    "express",
				Version: "4.18.2",
				Purl:    "pkg:npm/express@4.18.2",
				ExternalReferences: []cyclonedxextract.ExternalReference{
					{
						Type: "website",
						URL:  "", // Empty URL
					},
				},
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Should fall back to purl-generated URL
	expectedURL := "https://www.npmjs.com/package/express/v/4.18.2"
	if *attr.URL != expectedURL {
		t.Errorf("Expected URL to be purl-generated %q, got %q", expectedURL, *attr.URL)
	}
}

// TestExtractPackages_WithExternalRefNoPurl tests external ref without purl.
func TestExtractPackages_WithExternalRefNoPurl(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:    "custom-lib",
				Version: "1.0.0",
				Purl:    "",
				ExternalReferences: []cyclonedxextract.ExternalReference{
					{
						Type: "website",
						URL:  "https://example.com/custom-lib",
					},
				},
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Purl != "" {
		t.Errorf("Expected empty purl, got %q", attr.Purl)
	}

	if attr.URL == nil {
		t.Fatal("Expected URL to be set from external ref, got nil")
	}

	if *attr.URL != "https://example.com/custom-lib" {
		t.Errorf("Expected URL to be website ref 'https://example.com/custom-lib', got %q", *attr.URL)
	}
}

// TestExtractPackages_WithoutExternalRefs tests fallback to purl when no external refs.
func TestExtractPackages_WithoutExternalRefs(t *testing.T) {
	t.Parallel()

	bom := &cyclonedxextract.BOM{
		BOMFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Components: []cyclonedxextract.Component{
			{
				Name:               "lodash",
				Version:            "4.17.21",
				Purl:               "pkg:npm/lodash@4.17.21",
				ExternalReferences: []cyclonedxextract.ExternalReference{},
			},
		},
	}

	result := cyclonedxextract.ExtractPackages(bom)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Should fall back to purl-generated URL
	expectedURL := "https://www.npmjs.com/package/lodash/v/4.17.21"
	if *attr.URL != expectedURL {
		t.Errorf("Expected URL to be purl-generated %q, got %q", expectedURL, *attr.URL)
	}
}
