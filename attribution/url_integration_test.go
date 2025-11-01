//go:build integration

package attribution_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/boringbin/sbomattr/attribution"
)

// TestPurlToURL_HTTPValidation tests that generated URLs return valid HTTP responses.
// This integration test makes real HTTP requests to verify URLs work in practice.
func TestPurlToURL_HTTPValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		purl string
	}{
		// npm - popular package with org
		{
			name: "npm with org - types/node",
			purl: "pkg:npm/%40types/node@18.0.0",
		},
		// npm - popular package without org
		{
			name: "npm without org - lodash",
			purl: "pkg:npm/lodash@4.17.21",
		},
		// pypi - popular package
		{
			name: "pypi - requests",
			purl: "pkg:pypi/requests@2.31.0",
		},
		// cargo - popular Rust crate
		{
			name: "cargo - serde",
			purl: "pkg:cargo/serde@1.0.200",
		},
		// gem - popular Ruby gem
		{
			name: "gem - rails",
			purl: "pkg:gem/rails@7.0.0",
		},
		// golang - popular Go package
		{
			name: "golang - gin",
			purl: "pkg:golang/github.com/gin-gonic/gin@v1.9.0",
		},
		// nuget - popular .NET package
		{
			name: "nuget - Newtonsoft.Json",
			purl: "pkg:nuget/Newtonsoft.Json@13.0.1",
		},
		// pub - popular Dart package
		{
			name: "pub - http",
			purl: "pkg:pub/http@0.13.0",
		},
		// github - popular repository
		{
			name: "github - torvalds/linux",
			purl: "pkg:github/torvalds/linux@v6.0",
		},
		// composer - popular PHP package
		{
			name: "composer - symfony/symfony",
			purl: "pkg:composer/symfony/symfony@6.3.0",
		},
		// maven - popular Java package
		{
			name: "maven - junit",
			purl: "pkg:maven/junit/junit@4.13.2",
		},
		// docker - official nginx image
		{
			name: "docker - official nginx",
			purl: "pkg:docker/library/nginx@latest",
		},
		// docker - bitnami image
		{
			name: "docker - bitnami/nginx",
			purl: "pkg:docker/bitnami/nginx@latest",
		},
		// hex - popular Elixir package
		{
			name: "hex - phoenix",
			purl: "pkg:hex/phoenix@1.7.0",
		},
		// cocoapods - popular iOS library
		{
			name: "cocoapods - Alamofire",
			purl: "pkg:cocoapods/Alamofire@5.6.0",
		},
		// conda - popular data science package
		{
			name: "conda - numpy",
			purl: "pkg:conda/conda-forge/numpy@1.24.0",
		},
		// deb - common Debian package
		{
			name: "deb - curl",
			purl: "pkg:deb/debian/curl@7.88.1",
		},
		// rpm - common RPM package
		{
			name: "rpm - bash",
			purl: "pkg:rpm/redhat/bash@5.0",
		},
		// apk - common Alpine package
		{
			name: "apk - alpine",
			purl: "pkg:apk/alpine/alpine-base@3.18.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Generate the URL
			urlPtr := attribution.PurlToURL(tt.purl)
			if urlPtr == nil {
				t.Fatalf("PurlToURL returned nil for purl: %s", tt.purl)
			}

			url := *urlPtr

			// Verify the URL with HTTP request
			verifyURL(t, url, tt.purl)
		})
	}
}

// verifyURL makes an HTTP GET request to verify the URL returns a valid response.
// It accepts 2xx and 3xx status codes (including redirects).
// Network errors cause the test to skip (to avoid CI failures on network issues).
// 4xx and 5xx errors cause the test to fail (indicating wrong URL format).
func verifyURL(t *testing.T, url, purl string) {
	t.Helper()

	// Create context with timeout to avoid hanging tests
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request for %s: %v", url, err)
	}

	// Set user agent to identify the request
	req.Header.Set("User-Agent", "sbomattr-integration-test/1.0")

	// Make the request
	client := &http.Client{
		// Allow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		// Network errors should skip the test, not fail it
		// This prevents CI failures due to transient network issues
		t.Skipf("Network error accessing %s (purl: %s): %v", url, purl, err)
		return
	}
	defer resp.Body.Close()

	// Check status code
	// Accept 2xx (success) and 3xx (redirects) as valid
	// Fail on 4xx (client error - wrong URL) and 5xx (server error)
	if resp.StatusCode >= 400 {
		t.Errorf("URL %s returned error status %d %s (purl: %s)",
			url, resp.StatusCode, http.StatusText(resp.StatusCode), purl)
		return
	}

	// Log success for debugging
	t.Logf("✓ URL verified: %s (status: %d %s)", url, resp.StatusCode, http.StatusText(resp.StatusCode))
}

// TestPurlToURL_BitbucketHTTPValidation tests Bitbucket URLs separately
// as they may have different access patterns.
// Note: Many Bitbucket repositories have been migrated away or deleted,
// so this test may be skipped if the repository is not accessible.
func TestPurlToURL_BitbucketHTTPValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		purl string
	}{
		{
			name: "bitbucket - public repo",
			purl: "pkg:bitbucket/birkenfeld/pygments-main@2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			urlPtr := attribution.PurlToURL(tt.purl)
			if urlPtr == nil {
				t.Fatalf("PurlToURL returned nil for purl: %s", tt.purl)
			}

			url := *urlPtr

			// Bitbucket may require authentication for some repos,
			// and many repos have migrated away, so we're lenient here
			verifyURL(t, url, tt.purl)
		})
	}
}
