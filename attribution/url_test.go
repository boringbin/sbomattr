package attribution_test

import (
	"errors"
	"testing"

	"github.com/boringbin/sbomattr/attribution"
)

// TestPurlToURL_NPMWithOrg tests the PurlToURL function with an NPM package with an organization.
func TestPurlToURL_NPMWithOrg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		purl     string
		expected string
	}{
		{
			name:     "npm with org",
			purl:     "pkg:npm/%40babel/code-frame@7.22.5",
			expected: "https://www.npmjs.com/package/@babel/code-frame/v/7.22.5",
		},
		{
			name:     "npm without org",
			purl:     "pkg:npm/lodash@4.17.21",
			expected: "https://www.npmjs.com/package/lodash/v/4.17.21",
		},
		{
			name:     "npm with org - different package",
			purl:     "pkg:npm/%40types/node@18.0.0",
			expected: "https://www.npmjs.com/package/@types/node/v/18.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := attribution.PurlToURL(tt.purl, nil)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if result == nil {
				t.Fatalf("Expected URL, got nil")
			}

			if *result != tt.expected {
				t.Errorf("Expected URL %q, got %q", tt.expected, *result)
			}
		})
	}
}

// TestPurlToURL_OtherPackageTypes tests the PurlToURL function with other package types.
func TestPurlToURL_OtherPackageTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		purl     string
		expected string
	}{
		{
			name:     "cargo",
			purl:     "pkg:cargo/tokio@1.0.0",
			expected: "https://crates.io/crates/tokio/1.0.0",
		},
		{
			name:     "pypi",
			purl:     "pkg:pypi/django@4.2.0",
			expected: "https://pypi.org/project/django/4.2.0/",
		},
		{
			name:     "gem",
			purl:     "pkg:gem/rails@7.0.0",
			expected: "https://rubygems.org/gems/rails/versions/7.0.0",
		},
		{
			name:     "golang without namespace",
			purl:     "pkg:golang/github.com/gin-gonic/gin@v1.9.0",
			expected: "https://pkg.go.dev/github.com/gin-gonic/gin",
		},
		{
			name:     "nuget",
			purl:     "pkg:nuget/Newtonsoft.Json@13.0.1",
			expected: "https://www.nuget.org/packages/Newtonsoft.Json/13.0.1",
		},
		{
			name:     "pub",
			purl:     "pkg:pub/cookie_jar@4.0.8",
			expected: "https://pub.dev/packages/cookie_jar/versions/4.0.8",
		},
		{
			name:     "github with version tag",
			purl:     "pkg:github/golang/go@v1.21.0",
			expected: "https://github.com/golang/go/tree/v1.21.0",
		},
		{
			name:     "github with commit sha",
			purl:     "pkg:github/kubernetes/kubernetes@abc123def456",
			expected: "https://github.com/kubernetes/kubernetes/tree/abc123def456",
		},
		{
			name:     "composer",
			purl:     "pkg:composer/symfony/symfony@6.3.0",
			expected: "https://packagist.org/packages/symfony/symfony#6.3.0",
		},
		{
			name:     "maven",
			purl:     "pkg:maven/org.springframework/spring-core@5.3.28",
			expected: "https://central.sonatype.com/artifact/org.springframework/spring-core/5.3.28",
		},
		{
			name:     "golang with namespace",
			purl:     "pkg:golang/google.golang.org/grpc@v1.56.0",
			expected: "https://pkg.go.dev/google.golang.org/grpc",
		},
		{
			name:     "docker with namespace",
			purl:     "pkg:docker/bitnami/nginx@latest",
			expected: "https://hub.docker.com/r/bitnami/nginx",
		},
		{
			name:     "docker official image (library)",
			purl:     "pkg:docker/library/nginx@latest",
			expected: "https://hub.docker.com/_/nginx",
		},
		{
			name:     "docker without namespace",
			purl:     "pkg:docker/alpine@3.18",
			expected: "https://hub.docker.com/_/alpine",
		},
		{
			name:     "oci with namespace",
			purl:     "pkg:oci/bitnami/redis@7.0",
			expected: "https://hub.docker.com/r/bitnami/redis",
		},
		{
			name:     "oci official image",
			purl:     "pkg:oci/library/ubuntu@22.04",
			expected: "https://hub.docker.com/_/ubuntu",
		},
		{
			name:     "deb",
			purl:     "pkg:deb/debian/curl@7.88.1",
			expected: "https://packages.debian.org/curl",
		},
		{
			name:     "rpm",
			purl:     "pkg:rpm/redhat/openssl@1.1.1",
			expected: "https://rpmfind.net/linux/rpm2html/search.php?query=openssl",
		},
		{
			name:     "apk",
			purl:     "pkg:apk/alpine/curl@8.0.0",
			expected: "https://pkgs.alpinelinux.org/packages?name=curl",
		},
		{
			name:     "hex",
			purl:     "pkg:hex/phoenix@1.7.0",
			expected: "https://hex.pm/packages/phoenix/1.7.0",
		},
		{
			name:     "cocoapods",
			purl:     "pkg:cocoapods/Alamofire@5.6.0",
			expected: "https://cocoapods.org/pods/Alamofire",
		},
		{
			name:     "conda with namespace",
			purl:     "pkg:conda/conda-forge/numpy@1.24.0",
			expected: "https://anaconda.org/conda-forge/numpy",
		},
		{
			name:     "conda without namespace",
			purl:     "pkg:conda/pandas@2.0.0",
			expected: "https://anaconda.org/anaconda/pandas",
		},
		{
			name:     "bitbucket",
			purl:     "pkg:bitbucket/atlassian/python-bitbucket@0.1.0",
			expected: "https://bitbucket.org/atlassian/python-bitbucket/src/0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := attribution.PurlToURL(tt.purl, nil)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if result == nil {
				t.Fatalf("Expected URL, got nil")
			}

			if *result != tt.expected {
				t.Errorf("Expected URL %q, got %q", tt.expected, *result)
			}
		})
	}
}

// TestPurlToURL_InvalidPurl tests the PurlToURL function with an invalid purl.
func TestPurlToURL_InvalidPurl(t *testing.T) {
	t.Parallel()

	result, err := attribution.PurlToURL("not-a-valid-purl", nil)

	if err == nil {
		t.Error("Expected error for invalid purl, got nil")
	}

	// Should not be a sentinel error, but a parse error
	if errors.Is(err, attribution.ErrEmptyPurl) || errors.Is(err, attribution.ErrUnsupportedPurlType) {
		t.Errorf("Expected parse error, got sentinel error: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result for invalid purl, got %q", *result)
	}
}

// TestPurlToURL_EmptyPurl tests the PurlToURL function with an empty purl.
func TestPurlToURL_EmptyPurl(t *testing.T) {
	t.Parallel()

	result, err := attribution.PurlToURL("", nil)

	if !errors.Is(err, attribution.ErrEmptyPurl) {
		t.Errorf("Expected ErrEmptyPurl, got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil for empty purl, got %q", *result)
	}
}

// TestPurlToURL_UnsupportedType tests the PurlToURL function with an unsupported purl type.
func TestPurlToURL_UnsupportedType(t *testing.T) {
	t.Parallel()

	// Test various unsupported but recognized purl types
	unsupportedTypes := []struct {
		name string
		purl string
	}{
		{name: "alpm", purl: "pkg:alpm/arch/pacman@6.0.0"},
		{name: "bitnami", purl: "pkg:bitnami/nginx@1.0.0"},
		{name: "conan", purl: "pkg:conan/boost@1.76.0"},
		{name: "cran", purl: "pkg:cran/dplyr@1.0.0"},
		{name: "generic", purl: "pkg:generic/example@1.0.0"},
		{name: "hackage", purl: "pkg:hackage/aeson@2.0.0"},
		{name: "huggingface", purl: "pkg:huggingface/transformers@4.0.0"},
		{name: "mlflow", purl: "pkg:mlflow/model@1.0.0"},
	}

	for _, tt := range unsupportedTypes {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := attribution.PurlToURL(tt.purl, nil)

			if !errors.Is(err, attribution.ErrUnsupportedPurlType) {
				t.Errorf("Expected ErrUnsupportedPurlType for %q, got %v", tt.name, err)
			}

			if result != nil {
				t.Errorf("Expected nil for unsupported purl type %q, got %q", tt.name, *result)
			}
		})
	}
}

// TestPurlToURL_UnknownType tests the PurlToURL function with a completely unknown purl type.
func TestPurlToURL_UnknownType(t *testing.T) {
	t.Parallel()

	// Test with a completely unknown/made-up purl type
	result, err := attribution.PurlToURL("pkg:completely-unknown-type/package@1.0.0", nil)

	if !errors.Is(err, attribution.ErrUnsupportedPurlType) {
		t.Errorf("Expected ErrUnsupportedPurlType for unknown purl type, got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil for unknown purl type, got %q", *result)
	}
}
