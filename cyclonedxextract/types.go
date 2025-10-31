package cyclonedxextract

// See https://github.com/CycloneDX/cyclonedx-go

// BOM represents a minimal CycloneDX Bill of Materials with only the fields we need.
type BOM struct {
	BOMFormat   string      `json:"bomFormat"`
	SpecVersion string      `json:"specVersion"`
	Components  []Component `json:"components"`
}

// Component represents a minimal CycloneDX component with only the fields we need.
type Component struct {
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	Purl     string    `json:"purl"`
	Licenses *Licenses `json:"licenses"`
}

// Licenses represents the licenses field which can be structured in different ways.
type Licenses []LicenseChoice

// LicenseChoice represents a single license choice.
type LicenseChoice struct {
	License *License `json:"license"`
}

// License represents a license with various identification methods.
type License struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Expression string       `json:"expression"`
	Text       *LicenseText `json:"text"`
}

// LicenseText represents license text content.
type LicenseText struct {
	Content string `json:"content"`
}
