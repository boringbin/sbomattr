package spdxextract_test

import (
	"testing"

	"github.com/boringbin/sbomattr/spdxextract"
)

// TestExtractPackages_NilDocument tests the ExtractPackages function with a nil document.
func TestExtractPackages_NilDocument(t *testing.T) {
	t.Parallel()

	result := spdxextract.ExtractPackages(nil)

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(result))
	}
}

// TestExtractPackages_EmptyPackages tests the ExtractPackages function with an empty packages slice.
func TestExtractPackages_EmptyPackages(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages:    []spdxextract.Package{},
	}

	result := spdxextract.ExtractPackages(doc)

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(result))
	}
}

// TestExtractPackages_NilPackagesSlice tests the ExtractPackages function with a nil packages slice.
func TestExtractPackages_NilPackagesSlice(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages:    nil,
	}

	result := spdxextract.ExtractPackages(doc)

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(result))
	}
}

// TestExtractPackages_WithConcludedLicense tests the ExtractPackages function with a concluded license.
func TestExtractPackages_WithConcludedLicense(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "lodash",
				VersionInfo:      "4.17.21",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "Apache-2.0", // Should be ignored
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/lodash@4.17.21",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

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

// TestExtractPackages_WithDeclaredLicense_ConcludedIsNOASSERTION tests the ExtractPackages function with a declared license and concluded is NOASSERTION.
func TestExtractPackages_WithDeclaredLicense_ConcludedIsNOASSERTION(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "express",
				VersionInfo:      "4.18.2",
				LicenseConcluded: "NOASSERTION",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/express@4.18.2",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Name != "express" {
		t.Errorf("Expected name 'express', got %q", attr.Name)
	}

	if attr.License == nil {
		t.Fatal("Expected license to be set, got nil")
	}

	if *attr.License != "MIT" {
		t.Errorf("Expected license 'MIT' (declared), got %q", *attr.License)
	}
}

// TestExtractPackages_WithDeclaredLicense_ConcludedIsEmpty tests the ExtractPackages function with a declared license and concluded is empty.
func TestExtractPackages_WithDeclaredLicense_ConcludedIsEmpty(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "test-package",
				VersionInfo:      "1.0.0",
				LicenseConcluded: "",
				LicenseDeclared:  "Apache-2.0",
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.License == nil {
		t.Fatal("Expected license to be set, got nil")
	}

	if *attr.License != "Apache-2.0" {
		t.Errorf("Expected license 'Apache-2.0' (declared), got %q", *attr.License)
	}
}

// TestExtractPackages_WithPurlInExternalRefs tests the ExtractPackages function with a purl in external refs.
func TestExtractPackages_WithPurlInExternalRefs(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "react",
				VersionInfo:      "18.2.0",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/react@18.2.0",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Purl != "pkg:npm/react@18.2.0" {
		t.Errorf("Expected purl 'pkg:npm/react@18.2.0', got %q", attr.Purl)
	}
}

// TestExtractPackages_WithoutPurl tests the ExtractPackages function without a purl.
func TestExtractPackages_WithoutPurl(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "custom-package",
				VersionInfo:      "1.0.0",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs:     []spdxextract.ExternalRef{},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Purl != "" {
		t.Errorf("Expected empty purl, got %q", attr.Purl)
	}

	if attr.URL != nil {
		t.Errorf("Expected URL to be nil, got %q", *attr.URL)
	}
}

// TestExtractPackages_WithMultipleExternalRefs tests the ExtractPackages function with multiple external refs.
func TestExtractPackages_WithMultipleExternalRefs(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "multi-ref-package",
				VersionInfo:      "1.0.0",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "cpe23Type",
						ReferenceLocator: "cpe:2.3:a:vendor:product:1.0.0",
					},
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/multi-ref-package@1.0.0",
					},
					{
						ReferenceType:    "security-other",
						ReferenceLocator: "https://example.com/advisory",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	// Should extract the purl, not the other refs
	if attr.Purl != "pkg:npm/multi-ref-package@1.0.0" {
		t.Errorf("Expected purl 'pkg:npm/multi-ref-package@1.0.0', got %q", attr.Purl)
	}
}

// TestExtractPackages_NoPurlInExternalRefs tests the ExtractPackages function with no purl in external refs.
func TestExtractPackages_NoPurlInExternalRefs(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "no-purl-package",
				VersionInfo:      "1.0.0",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "cpe23Type",
						ReferenceLocator: "cpe:2.3:a:vendor:product:1.0.0",
					},
					{
						ReferenceType:    "security-other",
						ReferenceLocator: "https://example.com/advisory",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Purl != "" {
		t.Errorf("Expected empty purl, got %q", attr.Purl)
	}
}

// TestExtractPackages_MultiplePackages tests the ExtractPackages function with multiple packages.
func TestExtractPackages_MultiplePackages(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "package-one",
				VersionInfo:      "1.0.0",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/package-one@1.0.0",
					},
				},
			},
			{
				Name:             "package-two",
				VersionInfo:      "2.0.0",
				LicenseConcluded: "Apache-2.0",
				LicenseDeclared:  "Apache-2.0",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:pypi/package-two@2.0.0",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 2 {
		t.Fatalf("Expected 2 attributions, got %d", len(result))
	}

	if result[0].Name != "package-one" {
		t.Errorf("Expected first package name 'package-one', got %q", result[0].Name)
	}

	if result[1].Name != "package-two" {
		t.Errorf("Expected second package name 'package-two', got %q", result[1].Name)
	}

	if *result[0].License != "MIT" {
		t.Errorf("Expected first package license 'MIT', got %q", *result[0].License)
	}

	if *result[1].License != "Apache-2.0" {
		t.Errorf("Expected second package license 'Apache-2.0', got %q", *result[1].License)
	}
}

// TestExtractPackages_NilExternalRefs tests the ExtractPackages function with nil external refs.
func TestExtractPackages_NilExternalRefs(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "no-refs-package",
				VersionInfo:      "1.0.0",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs:     nil,
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Purl != "" {
		t.Errorf("Expected empty purl, got %q", attr.Purl)
	}

	if attr.URL != nil {
		t.Errorf("Expected URL to be nil, got %q", *attr.URL)
	}
}

// TestExtractPackages_WithHomepage tests that homepage is preferred over purl-generated URL.
func TestExtractPackages_WithHomepage(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "lodash",
				VersionInfo:      "4.17.21",
				Homepage:         "https://lodash.com",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/lodash@4.17.21",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Homepage should be preferred over purl-generated URL
	if *attr.URL != "https://lodash.com" {
		t.Errorf("Expected URL to be homepage 'https://lodash.com', got %q", *attr.URL)
	}
}

// TestExtractPackages_WithHomepageNONE tests that "NONE" homepage falls back to purl.
func TestExtractPackages_WithHomepageNONE(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "express",
				VersionInfo:      "4.18.2",
				Homepage:         "NONE",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/express@4.18.2",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

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

// TestExtractPackages_WithHomepageNOASSERTION tests that "NOASSERTION" homepage falls back to purl.
func TestExtractPackages_WithHomepageNOASSERTION(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "react",
				VersionInfo:      "18.2.0",
				Homepage:         "NOASSERTION",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/react@18.2.0",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Should fall back to purl-generated URL
	expectedURL := "https://www.npmjs.com/package/react/v/18.2.0"
	if *attr.URL != expectedURL {
		t.Errorf("Expected URL to be purl-generated %q, got %q", expectedURL, *attr.URL)
	}
}

// TestExtractPackages_WithHomepageEmpty tests that empty homepage falls back to purl.
func TestExtractPackages_WithHomepageEmpty(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "webpack",
				VersionInfo:      "5.88.2",
				Homepage:         "",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs: []spdxextract.ExternalRef{
					{
						ReferenceType:    "purl",
						ReferenceLocator: "pkg:npm/webpack@5.88.2",
					},
				},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.URL == nil {
		t.Fatal("Expected URL to be set, got nil")
	}

	// Should fall back to purl-generated URL
	expectedURL := "https://www.npmjs.com/package/webpack/v/5.88.2"
	if *attr.URL != expectedURL {
		t.Errorf("Expected URL to be purl-generated %q, got %q", expectedURL, *attr.URL)
	}
}

// TestExtractPackages_WithHomepageNoPurl tests homepage without purl.
func TestExtractPackages_WithHomepageNoPurl(t *testing.T) {
	t.Parallel()

	doc := &spdxextract.Document{
		SPDXVersion: "SPDX-2.3",
		SPDXID:      "SPDXRef-DOCUMENT",
		Packages: []spdxextract.Package{
			{
				Name:             "custom-lib",
				VersionInfo:      "1.0.0",
				Homepage:         "https://example.com/custom-lib",
				LicenseConcluded: "MIT",
				LicenseDeclared:  "MIT",
				ExternalRefs:     []spdxextract.ExternalRef{},
			},
		},
	}

	result := spdxextract.ExtractPackages(doc)

	if len(result) != 1 {
		t.Fatalf("Expected 1 attribution, got %d", len(result))
	}

	attr := result[0]
	if attr.Purl != "" {
		t.Errorf("Expected empty purl, got %q", attr.Purl)
	}

	if attr.URL == nil {
		t.Fatal("Expected URL to be set from homepage, got nil")
	}

	if *attr.URL != "https://example.com/custom-lib" {
		t.Errorf("Expected URL to be homepage 'https://example.com/custom-lib', got %q", *attr.URL)
	}
}
