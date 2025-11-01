package sbomattr_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/boringbin/sbomattr"
)

func TestProcess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid SPDX file",
			filename: "testdata/example-spdx.json",
			wantErr:  false,
		},
		{
			name:     "valid CycloneDX file",
			filename: "testdata/example-cyclonedx.json",
			wantErr:  false,
		},
		{
			name:     "GitHub-wrapped SPDX file",
			filename: "testdata/github-wrapped-spdx.json",
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(tc.filename)
			if err != nil {
				t.Fatalf("failed to read test file: %v", err)
			}

			ctx := context.Background()
			attrs, err := sbomattr.Process(ctx, data)

			if tc.wantErr && err == nil {
				t.Error("Process() expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Process() unexpected error: %v", err)
			}
			if !tc.wantErr && len(attrs) == 0 {
				t.Error("Process() returned empty attributions")
			}
		})
	}
}

func TestProcess_InvalidData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	invalidData := []byte(`{"invalid": "json"}`)

	_, err := sbomattr.Process(ctx, invalidData)
	if err == nil {
		t.Error("Process() with invalid data should return error")
	}
}

func TestProcess_Cancellation(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/example-spdx.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = sbomattr.Process(ctx, data)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Process() with cancelled context should return context.Canceled, got %v", err)
	}
}

func TestProcessFiles(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	filenames := []string{
		"testdata/example-spdx.json",
		"testdata/example-cyclonedx.json",
	}

	attrs, err := sbomattr.ProcessFiles(ctx, filenames)

	if err != nil {
		t.Errorf("ProcessFiles() unexpected error: %v", err)
	}
	if len(attrs) == 0 {
		t.Error("ProcessFiles() returned empty attributions")
	}
}

func TestProcessFiles_WithInvalidFiles(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	filenames := []string{
		"testdata/example-spdx.json",
		"testdata/does-not-exist.json", // This file doesn't exist
	}

	// Should still succeed because one valid file exists
	attrs, err := sbomattr.ProcessFiles(ctx, filenames)

	if err != nil {
		t.Errorf("ProcessFiles() unexpected error: %v", err)
	}
	if len(attrs) == 0 {
		t.Error("ProcessFiles() returned empty attributions despite valid file")
	}
}

// Integration test that processes all test files.
func TestProcessFiles_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	t.Parallel()

	ctx := context.Background()

	// Find all JSON files in testdata
	matches, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Fatalf("failed to glob testdata: %v", err)
	}

	if len(matches) == 0 {
		t.Skip("no test data files found")
	}

	attrs, err := sbomattr.ProcessFiles(ctx, matches)

	if err != nil {
		t.Errorf("ProcessFiles() unexpected error: %v", err)
	}
	if len(attrs) == 0 {
		t.Error("ProcessFiles() returned empty attributions")
	}

	t.Logf("Processed %d files and extracted %d deduplicated attributions", len(matches), len(attrs))
}

// TestSetLogger tests the SetLogger function with a valid logger.
func TestSetLogger(t *testing.T) {
	// Note: Cannot use t.Parallel() here because SetLogger modifies package-level state

	// Create a custom logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// This should not panic
	sbomattr.SetLogger(logger)

	// Test that operations still work with the new logger
	data, err := os.ReadFile("testdata/example-spdx.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	ctx := context.Background()
	attrs, err := sbomattr.Process(ctx, data)

	if err != nil {
		t.Errorf("Process() with custom logger unexpected error: %v", err)
	}

	if len(attrs) == 0 {
		t.Error("Process() with custom logger returned empty attributions")
	}
}

// TestSetLogger_Nil tests the SetLogger function with nil (should not crash).
func TestSetLogger_Nil(t *testing.T) {
	// Note: Cannot use t.Parallel() here because SetLogger modifies package-level state

	// This should not panic
	sbomattr.SetLogger(nil)

	// Test that operations still work after setting nil
	data, err := os.ReadFile("testdata/example-spdx.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	ctx := context.Background()
	attrs, err := sbomattr.Process(ctx, data)

	if err != nil {
		t.Errorf("Process() after SetLogger(nil) unexpected error: %v", err)
	}

	if len(attrs) == 0 {
		t.Error("Process() after SetLogger(nil) returned empty attributions")
	}
}

// TestProcess_InvalidSPDXJSON tests error handling when SPDX parsing fails.
func TestProcess_InvalidSPDXJSON(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	// This is truly malformed JSON that will fail parsing
	invalidSPDXData := []byte(`{"spdxVersion": "SPDX-2.3", this is broken`)

	_, err := sbomattr.Process(ctx, invalidSPDXData)
	if err == nil {
		t.Error("Process() with invalid SPDX JSON should return error")
	}
}

// TestProcess_InvalidCycloneDXJSON tests error handling when CycloneDX parsing fails.
func TestProcess_InvalidCycloneDXJSON(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	// This is truly malformed JSON that will fail parsing
	invalidCycloneDXData := []byte(`{"bomFormat": "CycloneDX", broken json here`)

	_, err := sbomattr.Process(ctx, invalidCycloneDXData)
	if err == nil {
		t.Error("Process() with invalid CycloneDX JSON should return error")
	}
}

// TestProcessFiles_AllInvalidFiles tests error when all files are invalid.
func TestProcessFiles_AllInvalidFiles(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	filenames := []string{
		"testdata/does-not-exist-1.json",
		"testdata/does-not-exist-2.json",
	}

	// Should return error because no valid attributions could be extracted
	attrs, err := sbomattr.ProcessFiles(ctx, filenames)

	if err == nil {
		t.Error("ProcessFiles() with all invalid files should return error")
	}

	if attrs != nil {
		t.Errorf("ProcessFiles() with all invalid files should return nil, got %+v", attrs)
	}
}

// TestProcessFiles_Cancellation tests context cancellation in ProcessFiles.
func TestProcessFiles_Cancellation(t *testing.T) {
	t.Parallel()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	filenames := []string{
		"testdata/example-spdx.json",
		"testdata/example-cyclonedx.json",
	}

	_, err := sbomattr.ProcessFiles(ctx, filenames)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("ProcessFiles() with cancelled context should return context.Canceled, got %v", err)
	}
}
