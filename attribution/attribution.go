package attribution

import (
	"io"
	"log/slog"
)

// logger is the package-level logger. By default, logging is disabled.
//
//nolint:gochecknoglobals // Package-level logger is simpler
var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

// SetLogger configures the logger for this package.
// By default, logging is disabled. Call this function to enable logging.
func SetLogger(l *slog.Logger) {
	if l != nil {
		logger = l
	}
}

// Attribution represents a simplified view of an SBOM package with essential fields.
//
// The goal is to be able to use this to point to:
// - Describe the package
// - Outline it's license and usage restrictions
// - Provide a way to confirm the information yourself.
type Attribution struct {
	// Name is the package name
	Name string `json:"name"`
	// License is the declared license
	License *string `json:"license,omitempty"`
	// URL is the package URL
	URL *string `json:"url,omitempty"`
	// Purl is the package purl
	Purl string `json:"purl"`
}
