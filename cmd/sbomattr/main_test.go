package main

import (
	"bytes"
	"flag"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPrintUsage tests the printUsage function.
func TestPrintUsage(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	printUsage(&buf, "sbomattr")
	output := buf.String()

	// Check that usage contains expected strings
	expectedStrings := []string{
		"Usage:",
		"Create an aggregated notice",
		"Arguments:",
		"file-or-directory",
		"Options:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("printUsage() output missing %q\nGot output:\n%s", expected, output)
		}
	}
}

// TestSetupLogger tests the setupLogger function.
func TestSetupLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		verbose bool
		want    slog.Level
	}{
		{
			name:    "verbose mode",
			verbose: true,
			want:    slog.LevelDebug,
		},
		{
			name:    "non-verbose mode",
			verbose: false,
			want:    slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := setupLogger(tt.verbose)
			if logger == nil {
				t.Fatal("setupLogger() returned nil")
			}

			// Logger should be configured but we can't easily inspect the level
			// We mainly test that it doesn't panic and returns a logger
		})
	}
}

// TestExpandPaths_SingleFile tests expandPaths with a single file.
func TestExpandPaths_SingleFile(t *testing.T) {
	t.Parallel()

	// Create a temporary file
	tmpFile, err := os.CreateTemp(t.TempDir(), "sbom-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	logger := setupLogger(false)
	files := expandPaths([]string{tmpFile.Name()}, logger)

	if len(files) != 1 {
		t.Errorf("expandPaths() returned %d files, want 1", len(files))
	}

	if len(files) > 0 && files[0] != tmpFile.Name() {
		t.Errorf("expandPaths() = %v, want %v", files[0], tmpFile.Name())
	}
}

// TestExpandPaths_Directory tests expandPaths with a directory containing JSON files.
func TestExpandPaths_Directory(t *testing.T) {
	t.Parallel()

	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create some test files
	jsonFile1 := filepath.Join(tmpDir, "test1.json")
	jsonFile2 := filepath.Join(tmpDir, "test2.json")
	txtFile := filepath.Join(tmpDir, "test.txt")

	for _, file := range []string{jsonFile1, jsonFile2, txtFile} {
		if createErr := os.WriteFile(file, []byte("{}"), 0600); createErr != nil {
			t.Fatalf("failed to create test file: %v", createErr)
		}
	}

	logger := setupLogger(false)
	files := expandPaths([]string{tmpDir}, logger)

	// Should only include .json files
	expectedCount := 2
	if len(files) != expectedCount {
		t.Errorf("expandPaths() returned %d files, want %d", len(files), expectedCount)
	}

	// Check that both JSON files are included
	foundFiles := make(map[string]bool)
	for _, f := range files {
		foundFiles[filepath.Base(f)] = true
	}

	if !foundFiles["test1.json"] || !foundFiles["test2.json"] {
		t.Errorf("expandPaths() = %v, want test1.json and test2.json", files)
	}

	if foundFiles["test.txt"] {
		t.Error("expandPaths() should not include .txt files")
	}
}

// TestExpandPaths_NonExistentPath tests expandPaths with non-existent path.
func TestExpandPaths_NonExistentPath(t *testing.T) {
	t.Parallel()

	logger := setupLogger(false)
	files := expandPaths([]string{"/nonexistent/path/to/file.json"}, logger)

	// Should return empty slice for non-existent paths
	if len(files) != 0 {
		t.Errorf("expandPaths() with non-existent path returned %d files, want 0", len(files))
	}
}

// TestExpandPaths_EmptyDirectory tests expandPaths with an empty directory.
func TestExpandPaths_EmptyDirectory(t *testing.T) {
	t.Parallel()

	// Create an empty temporary directory
	tmpDir := t.TempDir()

	logger := setupLogger(false)
	files := expandPaths([]string{tmpDir}, logger)

	if len(files) != 0 {
		t.Errorf("expandPaths() with empty directory returned %d files, want 0", len(files))
	}
}

// TestExpandPaths_MixedPaths tests expandPaths with mixed files and directories.
func TestExpandPaths_MixedPaths(t *testing.T) {
	t.Parallel()

	// Create temporary directory
	tmpDir := t.TempDir()

	// Create a JSON file in the directory
	dirFile := filepath.Join(tmpDir, "dir-file.json")
	if createErr := os.WriteFile(dirFile, []byte("{}"), 0600); createErr != nil {
		t.Fatalf("failed to create dir file: %v", createErr)
	}

	// Create a standalone JSON file
	tmpFile, err := os.CreateTemp(t.TempDir(), "sbom-standalone-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	logger := setupLogger(false)
	files := expandPaths([]string{tmpDir, tmpFile.Name()}, logger)

	// Should return both the file from directory and the standalone file
	expectedCount := 2
	if len(files) != expectedCount {
		t.Errorf("expandPaths() returned %d files, want %d", len(files), expectedCount)
	}
}

// TestRun_Version tests the run function with the --version flag.
func TestRun_Version(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"sbomattr", "--version"}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := run()

	_ = w.Close()
	os.Stdout = oldStdout

	if exitCode != exitSuccess {
		t.Errorf("run() with --version returned exit code %d, want %d", exitCode, exitSuccess)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "sbomattr version") {
		t.Errorf("run() --version output = %q, want to contain 'sbomattr version'", output)
	}
	if !strings.Contains(output, version) {
		t.Errorf("run() --version output = %q, want to contain version %q", output, version)
	}
}

// TestRun_NoArguments tests the run function with no arguments.
func TestRun_NoArguments(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"sbomattr"}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	exitCode := run()

	_ = w.Close()
	os.Stderr = oldStderr

	if exitCode != exitInvalidArgs {
		t.Errorf("run() with no args returned exit code %d, want %d", exitCode, exitInvalidArgs)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "no SBOM files") {
		t.Errorf("run() no args stderr should mention no SBOM files, got: %s", output)
	}
}

// TestRun_ValidSingleFile tests the run function with a single valid SBOM file.
func TestRun_ValidSingleFile(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Use the existing test data
	testFile := "../../testdata/example-spdx.json"
	os.Args = []string{"sbomattr", testFile}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := run()

	_ = w.Close()
	os.Stdout = oldStdout

	if exitCode != exitSuccess {
		t.Errorf("run() with valid SBOM returned exit code %d, want %d", exitCode, exitSuccess)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Check for CSV header
	if !strings.Contains(output, "Name,License,Purl,URL") {
		t.Errorf("run() output should contain CSV header, got: %s", output)
	}

	// Check for at least one package name from the SPDX file
	if !strings.Contains(output, "lodash") && !strings.Contains(output, "react") {
		t.Errorf("run() output should contain package names, got: %s", output)
	}
}

// TestRun_ValidMultipleFiles tests the run function with multiple valid SBOM files.
func TestRun_ValidMultipleFiles(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Use multiple test files
	testFile1 := "../../testdata/example-spdx.json"
	testFile2 := "../../testdata/example-cyclonedx.json"
	os.Args = []string{"sbomattr", testFile1, testFile2}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := run()

	_ = w.Close()
	os.Stdout = oldStdout

	if exitCode != exitSuccess {
		t.Errorf("run() with multiple valid SBOMs returned exit code %d, want %d", exitCode, exitSuccess)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Check for CSV header
	if !strings.Contains(output, "Name,License,Purl,URL") {
		t.Errorf("run() output should contain CSV header")
	}

	// Should contain packages from both files (deduplication may occur)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	// At least header + some data rows
	if len(lines) < 2 {
		t.Errorf("run() output should contain multiple rows, got %d", len(lines))
	}
}

// TestRun_ValidDirectory tests the run function with a directory.
func TestRun_ValidDirectory(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Use the testdata directory
	testDir := "../../testdata"
	os.Args = []string{"sbomattr", testDir}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := run()

	_ = w.Close()
	os.Stdout = oldStdout

	if exitCode != exitSuccess {
		t.Errorf("run() with valid directory returned exit code %d, want %d", exitCode, exitSuccess)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Check for CSV header
	if !strings.Contains(output, "Name,License,Purl,URL") {
		t.Errorf("run() output should contain CSV header")
	}
}

// TestRun_InvalidSBOM tests the run function with an invalid SBOM file.
func TestRun_InvalidSBOM(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Create a temporary file with invalid JSON
	tmpFile, err := os.CreateTemp(t.TempDir(), "invalid-sbom-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write invalid JSON
	_, _ = tmpFile.WriteString("{this is not valid json")
	tmpFile.Close()

	os.Args = []string{"sbomattr", tmpFile.Name()}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	exitCode := run()

	_ = w.Close()
	os.Stderr = oldStderr

	if exitCode != exitInvalidSBOM {
		t.Errorf("run() with invalid SBOM returned exit code %d, want %d", exitCode, exitInvalidSBOM)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "failed to process SBOM") {
		t.Errorf("run() stderr should mention failed to process SBOM, got: %s", output)
	}
}

// TestRun_NonExistentFile tests the run function with a non-existent file.
func TestRun_NonExistentFile(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"sbomattr", "/nonexistent/file.json"}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	exitCode := run()

	_ = w.Close()
	os.Stderr = oldStderr

	if exitCode != exitInvalidArgs {
		t.Errorf("run() with non-existent file returned exit code %d, want %d", exitCode, exitInvalidArgs)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Should log error about not being able to access the path
	if !strings.Contains(output, "cannot access path") && !strings.Contains(output, "no SBOM files found") {
		t.Errorf("run() stderr should mention path access error, got: %s", output)
	}
}

// TestRun_VerboseMode tests the run function with verbose flag.
func TestRun_VerboseMode(t *testing.T) {
	// Note: Cannot use t.Parallel() because run() modifies global flag.CommandLine

	// Save and restore os.Args and flag.CommandLine
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	testFile := "../../testdata/example-spdx.json"
	os.Args = []string{"sbomattr", "-v", testFile}

	// Capture stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	exitCode := run()

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	if exitCode != exitSuccess {
		t.Errorf("run() with -v flag returned exit code %d, want %d", exitCode, exitSuccess)
	}

	var bufOut bytes.Buffer
	var bufErr bytes.Buffer
	_, _ = io.Copy(&bufOut, rOut)
	_, _ = io.Copy(&bufErr, rErr)

	// Stdout should contain CSV output
	if !strings.Contains(bufOut.String(), "Name,License,Purl,URL") {
		t.Error("run() stdout should contain CSV output in verbose mode")
	}

	// Stderr may contain debug logs (depending on logger configuration)
	// We just verify the command runs successfully
}

// TestExpandPaths_DirectoryWithSubdirectories tests that subdirectories are not recursively searched.
func TestExpandPaths_DirectoryWithSubdirectories(t *testing.T) {
	t.Parallel()

	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create a JSON file in the root directory
	rootFile := filepath.Join(tmpDir, "root.json")
	if createErr := os.WriteFile(rootFile, []byte("{}"), 0600); createErr != nil {
		t.Fatalf("failed to create root file: %v", createErr)
	}

	// Create a subdirectory with a JSON file
	subDir := filepath.Join(tmpDir, "subdir")
	if mkdirErr := os.Mkdir(subDir, 0700); mkdirErr != nil {
		t.Fatalf("failed to create subdir: %v", mkdirErr)
	}

	subFile := filepath.Join(subDir, "sub.json")
	if createErr := os.WriteFile(subFile, []byte("{}"), 0600); createErr != nil {
		t.Fatalf("failed to create sub file: %v", createErr)
	}

	logger := setupLogger(false)
	files := expandPaths([]string{tmpDir}, logger)

	// Should only include root.json, not sub.json (non-recursive)
	if len(files) != 1 {
		t.Errorf("expandPaths() returned %d files, want 1 (non-recursive)", len(files))
	}

	if len(files) > 0 && filepath.Base(files[0]) != "root.json" {
		t.Errorf("expandPaths() = %v, want root.json", filepath.Base(files[0]))
	}
}
