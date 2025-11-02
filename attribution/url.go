package attribution

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/package-url/packageurl-go"
)

// Sentinel errors for PurlToURL function.
var (
	// ErrEmptyPurl is returned when the purl string is empty or whitespace-only.
	ErrEmptyPurl = errors.New("empty purl string")
	// ErrUnsupportedPurlType is returned when the purl type is not supported for URL generation.
	ErrUnsupportedPurlType = errors.New("unsupported purl type")
)

// PurlToURL constructs a package management URL from a purl string.
// Returns ErrEmptyPurl if the purl string is empty or whitespace-only.
// Returns ErrUnsupportedPurlType if the purl type is not supported for URL generation.
// Returns other errors if the purl string is malformed.
// The logger parameter is optional; pass nil to disable logging.
func PurlToURL(purlString string, logger *slog.Logger) (*string, error) {
	if strings.TrimSpace(purlString) == "" {
		return nil, ErrEmptyPurl
	}

	purl, err := packageurl.FromString(purlString)
	if err != nil {
		return nil, fmt.Errorf("parse purl: %w", err)
	}

	return mapPurlToURL(purl, logger)
}

// mapPurlToURL maps a purl to a package management URL.
func mapPurlToURL(purl packageurl.PackageURL, logger *slog.Logger) (*string, error) {
	// See https://github.com/package-url/purl-spec#known-purl-types
	switch purl.Type {
	case "cargo":
		return buildURL("https://crates.io/crates/%s/%s", purl.Name, purl.Version), nil
	case "composer":
		return buildURL("https://packagist.org/packages/%s/%s#%s", purl.Namespace, purl.Name, purl.Version), nil
	case "gem":
		return buildURL("https://rubygems.org/gems/%s/versions/%s", purl.Name, purl.Version), nil
	case "golang":
		if purl.Namespace != "" {
			return buildURL("https://pkg.go.dev/%s/%s@%s", purl.Namespace, purl.Name, purl.Version), nil
		}
		return buildURL("https://pkg.go.dev/%s@%s", purl.Name, purl.Version), nil
	case "maven":
		return buildURL("https://mvnrepository.com/artifact/%s/%s/%s", purl.Namespace, purl.Name, purl.Version), nil
	case "npm":
		if purl.Namespace != "" {
			return buildURL("https://www.npmjs.com/package/%s/%s/v/%s", purl.Namespace, purl.Name, purl.Version), nil
		}
		return buildURL("https://www.npmjs.com/package/%s/v/%s", purl.Name, purl.Version), nil
	case "nuget":
		return buildURL("https://www.nuget.org/packages/%s/%s", purl.Name, purl.Version), nil
	case "pub":
		return buildURL("https://pub.dev/packages/%s/versions/%s", purl.Name, purl.Version), nil
	case "pypi":
		return buildURL("https://pypi.org/project/%s/%s/", purl.Name, purl.Version), nil
	case "github":
		return buildURL("https://github.com/%s/%s/tree/%s", purl.Namespace, purl.Name, purl.Version), nil
	case "docker", "oci":
		return buildDockerHubURL(purl.Namespace, purl.Name), nil
	case "deb":
		return buildURL("https://packages.debian.org/%s", purl.Name), nil
	case "rpm":
		return buildURL("https://rpmfind.net/linux/rpm2html/search.php?query=%s", purl.Name), nil
	case "apk":
		return buildURL("https://pkgs.alpinelinux.org/packages?name=%s", purl.Name), nil
	case "hex":
		return buildURL("https://hex.pm/packages/%s/%s", purl.Name, purl.Version), nil
	case "cocoapods":
		return buildURL("https://cocoapods.org/pods/%s", purl.Name), nil
	case "conda":
		if purl.Namespace != "" {
			return buildURL("https://anaconda.org/%s/%s", purl.Namespace, purl.Name), nil
		}
		return buildURL("https://anaconda.org/anaconda/%s", purl.Name), nil
	case "bitbucket":
		return buildURL("https://bitbucket.org/%s/%s/src/%s", purl.Namespace, purl.Name, purl.Version), nil
	default:
		if logger != nil {
			logger.Debug("purl type not supported", "type", purl.Type)
		}
		return nil, ErrUnsupportedPurlType
	}
}

// buildURL constructs a URL from a format string and arguments.
func buildURL(format string, args ...any) *string {
	url := fmt.Sprintf(format, args...)
	return &url
}

// buildDockerHubURL constructs a Docker Hub URL for docker/oci images.
// Official images (library namespace) use the "_" prefix, others use "r/" prefix.
func buildDockerHubURL(namespace, name string) *string {
	if namespace != "" && namespace != "library" {
		return buildURL("https://hub.docker.com/r/%s/%s", namespace, name)
	}
	return buildURL("https://hub.docker.com/_/%s", name)
}
