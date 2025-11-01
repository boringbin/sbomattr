package spdxextract

// See https://github.com/spdx/tools-golang

// Document represents a minimal SPDX document with only the fields we need.
type Document struct {
	SPDXVersion string    `json:"spdxVersion"`
	SPDXID      string    `json:"SPDXID"`
	Packages    []Package `json:"packages"`
}

// Package represents a minimal SPDX package with only the fields we need.
type Package struct {
	Name             string        `json:"name"`
	VersionInfo      string        `json:"versionInfo"`
	Homepage         string        `json:"homepage"`
	LicenseConcluded string        `json:"licenseConcluded"`
	LicenseDeclared  string        `json:"licenseDeclared"`
	ExternalRefs     []ExternalRef `json:"externalRefs"`
}

// ExternalRef represents an external reference (like purl).
type ExternalRef struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceType     string `json:"referenceType"`
	ReferenceLocator  string `json:"referenceLocator"`
}
