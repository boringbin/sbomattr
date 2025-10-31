package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/boringbin/sbomattr"
	"github.com/boringbin/sbomattr/attribution"
	"github.com/boringbin/sbomattr/format"
)

const (
	// version is the version of the program.
	version = "0.1.0-dev"
	// exitSuccess is the exit code for success.
	exitSuccess = 0
	// exitInvalidArgs is the exit code for invalid arguments.
	exitInvalidArgs = 1
	// exitInvalidSBOM is the exit code for invalid SBOM.
	exitInvalidSBOM = 2
	// exitRuntimeError is the exit code for runtime error.
	exitRuntimeError = 3
)

func main() {
	os.Exit(run())
}

func run() int {
	var (
		verbose     = flag.Bool("v", false, "Verbose output (debug mode)")
		showVersion = flag.Bool("version", false, "Show version and exit")
	)

	// Customize usage message
	printUsageFunc := func() {
		printUsage()
	}
	flag.CommandLine.Usage = printUsageFunc

	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Fprintf(os.Stdout, "sbomattr version %s\n", version)
		return exitSuccess
	}

	// Setup logger based on verbose flag
	logger := setupLogger(*verbose)

	// Get the input paths from the arguments
	args := flag.Args()

	// Validate arguments
	if len(args) == 0 {
		logger.Error("no SBOM files or directories provided")
		printUsage()
		return exitInvalidArgs
	}

	// Expand paths to get list of files
	files := expandPaths(args, logger)

	if len(files) == 0 {
		logger.Error("no SBOM files found")
		return exitInvalidArgs
	}

	// Configure package-level loggers
	sbomattr.SetLogger(logger)
	attribution.SetLogger(logger)

	// Process all files using the library
	ctx := context.Background()
	attributions, err := sbomattr.ProcessFiles(ctx, files)
	if err != nil {
		logger.Error("failed to process SBOM files", "error", err)
		return exitInvalidSBOM
	}

	// Output as CSV
	err = format.CSV(os.Stdout, attributions)
	if err != nil {
		logger.Error("failed to write CSV output", "error", err)
		return exitRuntimeError
	}

	return exitSuccess
}

// printUsage prints the usage message.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <file-or-directory>...\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Create an aggregated notice for one or more SBOMs.\n\n")
	fmt.Fprintf(os.Stderr, "Arguments:\n")
	fmt.Fprintf(
		os.Stderr,
		"  file-or-directory   SBOM files or directories containing SBOM files\n\n",
	)
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

// setupLogger sets up the logger based on the verbose flag.
func setupLogger(verbose bool) *slog.Logger {
	logLevel := slog.LevelError
	if verbose {
		// If verbose is true, set the log level to debug
		// This will log all messages, including debug messages
		logLevel = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
}

// expandPaths takes a mix of files and directories and returns a list of SBOM file paths.
func expandPaths(paths []string, logger *slog.Logger) []string {
	var files []string

	for _, path := range paths {
		info, statErr := os.Stat(path)
		if statErr != nil {
			logger.Error("cannot access path", "path", path, "error", statErr)
			continue
		}

		if info.IsDir() {
			// Read directory (non-recursive)
			entries, readErr := os.ReadDir(path)
			if readErr != nil {
				logger.Error("cannot read directory", "path", path, "error", readErr)
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				// Only consider JSON files (SBOM files are typically JSON)
				if strings.HasSuffix(entry.Name(), ".json") {
					files = append(files, filepath.Join(path, entry.Name()))
				}
			}
		} else {
			// Regular file
			files = append(files, path)
		}
	}

	return files
}
