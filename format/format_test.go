package format_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/boringbin/sbomattr/attribution"
	"github.com/boringbin/sbomattr/format"
)

// TestCSV tests the CSV function.
func TestCSV(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input []attribution.Attribution
		want  string
	}{
		{
			name:  "empty slice",
			input: []attribution.Attribution{},
			want:  "Name,License,Purl,URL\n",
		},
		{
			name: "single attribution with all fields",
			input: []attribution.Attribution{
				{
					Name:    "test-package",
					License: strPtr("MIT"),
					Purl:    "pkg:npm/test-package@1.0.0",
					URL:     strPtr("https://www.npmjs.com/package/test-package"),
				},
			},
			want: "Name,License,Purl,URL\n" +
				"test-package,MIT,pkg:npm/test-package@1.0.0,https://www.npmjs.com/package/test-package\n",
		},
		{
			name: "attribution with nil license and URL",
			input: []attribution.Attribution{
				{
					Name:    "test-package",
					License: nil,
					Purl:    "pkg:npm/test-package@1.0.0",
					URL:     nil,
				},
			},
			want: "Name,License,Purl,URL\n" +
				"test-package,,pkg:npm/test-package@1.0.0,\n",
		},
		{
			name: "attribution with commas in name",
			input: []attribution.Attribution{
				{
					Name:    "package, with, commas",
					License: strPtr("MIT"),
					Purl:    "pkg:npm/package-with-commas@1.0.0",
					URL:     nil,
				},
			},
			want: "Name,License,Purl,URL\n" +
				"\"package, with, commas\",MIT,pkg:npm/package-with-commas@1.0.0,\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := format.CSV(&buf, tc.input)
			if err != nil {
				t.Fatalf("CSV() unexpected error: %v", err)
			}

			if buf.String() != tc.want {
				t.Errorf("CSV() = %q, want %q", buf.String(), tc.want)
			}
		})
	}
}

// TestJSON tests the JSON function.
func TestJSON(t *testing.T) {
	t.Parallel()

	input := []attribution.Attribution{
		{
			Name:    "test-package",
			License: strPtr("MIT"),
			Purl:    "pkg:npm/test-package@1.0.0",
			URL:     strPtr("https://example.com/test-package"),
		},
	}

	var buf bytes.Buffer
	err := format.JSON(&buf, input)
	if err != nil {
		t.Fatalf("JSON() unexpected error: %v", err)
	}

	// Verify it contains expected fields
	output := buf.String()
	if !strings.Contains(output, "test-package") {
		t.Error("JSON() output should contain package name")
	}
	if !strings.Contains(output, "MIT") {
		t.Error("JSON() output should contain license")
	}
	if !strings.Contains(output, "pkg:npm/test-package@1.0.0") {
		t.Error("JSON() output should contain purl")
	}
}

// failingWriter is a mock writer that always returns an error.
type failingWriter struct{}

func (w *failingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("mock write error")
}

// TestJSON_WriteError tests JSON error handling when writer fails.
func TestJSON_WriteError(t *testing.T) {
	t.Parallel()

	input := []attribution.Attribution{
		{
			Name:    "test-package",
			License: strPtr("MIT"),
			Purl:    "pkg:npm/test-package@1.0.0",
			URL:     strPtr("https://example.com"),
		},
	}

	// Writer that fails immediately
	writer := &failingWriter{}
	err := format.JSON(writer, input)

	if err == nil {
		t.Error("JSON() with failing writer should return error")
	}
}

// strPtr converts a string to a pointer to a string.
func strPtr(s string) *string {
	return &s
}
