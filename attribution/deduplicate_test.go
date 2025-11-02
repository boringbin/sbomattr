package attribution_test

import (
	"testing"

	"github.com/boringbin/sbomattr/attribution"
)

// TestDeduplicate tests the Deduplicate function.
func TestDeduplicate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input []attribution.Attribution
		want  []attribution.Attribution
	}{
		{
			name:  "empty slice",
			input: []attribution.Attribution{},
			want:  []attribution.Attribution{},
		},
		{
			name: "no duplicates",
			input: []attribution.Attribution{
				{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg2", Purl: "pkg:npm/pkg2@2.0.0"},
			},
			want: []attribution.Attribution{
				{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg2", Purl: "pkg:npm/pkg2@2.0.0"},
			},
		},
		{
			name: "duplicates by purl",
			input: []attribution.Attribution{
				{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg1-duplicate", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg2", Purl: "pkg:npm/pkg2@2.0.0"},
			},
			want: []attribution.Attribution{
				{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg2", Purl: "pkg:npm/pkg2@2.0.0"},
			},
		},
		{
			name: "duplicates by name (no purl)",
			input: []attribution.Attribution{
				{Name: "pkg1", Purl: ""},
				{Name: "pkg1", Purl: ""},
				{Name: "pkg2", Purl: ""},
			},
			want: []attribution.Attribution{
				{Name: "pkg1", Purl: ""},
				{Name: "pkg2", Purl: ""},
			},
		},
		{
			name: "mixed purl and name keys",
			input: []attribution.Attribution{
				{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg2", Purl: ""},
				{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg2", Purl: ""},
			},
			want: []attribution.Attribution{
				{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
				{Name: "pkg2", Purl: ""},
			},
		},
		{
			name: "preserves first occurrence",
			input: []attribution.Attribution{
				{Name: "first", Purl: "pkg:npm/pkg@1.0.0", License: strPtr("MIT")},
				{Name: "second", Purl: "pkg:npm/pkg@1.0.0", License: strPtr("Apache-2.0")},
			},
			want: []attribution.Attribution{
				{Name: "first", Purl: "pkg:npm/pkg@1.0.0", License: strPtr("MIT")},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := attribution.Deduplicate(tc.input, nil)

			if len(got) != len(tc.want) {
				t.Errorf("Deduplicate() length = %d, want %d", len(got), len(tc.want))
			}

			for i := range got {
				if i >= len(tc.want) {
					break
				}
				if got[i].Name != tc.want[i].Name {
					t.Errorf("Deduplicate()[%d].Name = %q, want %q", i, got[i].Name, tc.want[i].Name)
				}
				if got[i].Purl != tc.want[i].Purl {
					t.Errorf("Deduplicate()[%d].Purl = %q, want %q", i, got[i].Purl, tc.want[i].Purl)
				}
			}
		})
	}
}

// TestDeduplicate_NilLogger tests the Deduplicate function works correctly with nil logger.
func TestDeduplicate_NilLogger(t *testing.T) {
	t.Parallel()

	input := []attribution.Attribution{
		{Name: "pkg1", Purl: "pkg:npm/pkg1@1.0.0"},
		{Name: "pkg1-duplicate", Purl: "pkg:npm/pkg1@1.0.0"},
	}

	got := attribution.Deduplicate(input, nil)

	const expectedLength = 1
	if len(got) != expectedLength {
		t.Errorf("Deduplicate() length = %d, want %d", len(got), expectedLength)
	}
}

// strPtr converts a string to a pointer to a string.
func strPtr(s string) *string {
	return &s
}
