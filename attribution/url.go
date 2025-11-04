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
		return buildCargoURL(purl), nil
	case "composer":
		return buildComposerURL(purl), nil
	case "gem":
		return buildGemURL(purl), nil
	case "golang":
		return buildGolangURL(purl), nil
	case "maven":
		return buildMavenURL(purl), nil
	case "npm":
		return buildNPMURL(purl), nil
	case "nuget":
		return buildNugetURL(purl), nil
	case "pub":
		return buildPubURL(purl), nil
	case "pypi":
		return buildPypiURL(purl), nil
	case "github":
		return buildGithubURL(purl), nil
	case "docker", "oci":
		return buildDockerHubURL(purl), nil
	case "deb":
		return buildDebURL(purl), nil
	case "rpm":
		return buildRpmURL(purl), nil
	case "apk":
		return buildApkURL(purl), nil
	case "hex":
		return buildHexURL(purl), nil
	case "cocoapods":
		return buildCocoapodsURL(purl), nil
	case "conda":
		return buildCondaURL(purl), nil
	case "bitbucket":
		return buildBitbucketURL(purl), nil
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

// buildCargoURL constructs a Cargo package URL from a purl.
// https://crates.io allows you to specify a version in the URL following the package name.
func buildCargoURL(purl packageurl.PackageURL) *string {
	return buildURL("https://crates.io/crates/%s/%s", purl.Name, purl.Version)
}

// buildComposerURL constructs a Composer package URL from a purl.
// https://packagist.org allows you to select a version and it will appear as an anchor in the URL.
func buildComposerURL(purl packageurl.PackageURL) *string {
	return buildURL("https://packagist.org/packages/%s/%s#%s", purl.Namespace, purl.Name, purl.Version)
}

// buildGemURL constructs a RubyGems package URL from a purl.
func buildGemURL(purl packageurl.PackageURL) *string {
	return buildURL("https://rubygems.org/gems/%s/versions/%s", purl.Name, purl.Version)
}

// buildGolangURL constructs a Go package URL from a purl.
// Version is not used, since versions are constructed in the https://pkg.go.dev documentation using tags.
// Most packages use a prefix like `v1.0.0`, but this isn't always the case.
func buildGolangURL(purl packageurl.PackageURL) *string {
	if purl.Namespace != "" {
		return buildURL("https://pkg.go.dev/%s/%s", purl.Namespace, purl.Name)
	}
	return buildURL("https://pkg.go.dev/%s", purl.Name)
}

// buildMavenURL constructs a Maven package URL from a purl.
// Uses the Maven Central repository URL.
func buildMavenURL(purl packageurl.PackageURL) *string {
	return buildURL("https://central.sonatype.com/artifact/%s/%s/%s", purl.Namespace, purl.Name, purl.Version)
}

// buildNPMURL constructs an NPM package URL from a purl.
func buildNPMURL(purl packageurl.PackageURL) *string {
	if purl.Namespace != "" {
		return buildURL("https://www.npmjs.com/package/%s/%s/v/%s", purl.Namespace, purl.Name, purl.Version)
	}
	return buildURL("https://www.npmjs.com/package/%s/v/%s", purl.Name, purl.Version)
}

// buildNugetURL constructs a NuGet package URL from a purl.
func buildNugetURL(purl packageurl.PackageURL) *string {
	return buildURL("https://www.nuget.org/packages/%s/%s", purl.Name, purl.Version)
}

// buildPubURL constructs a Pub package URL from a purl.
func buildPubURL(purl packageurl.PackageURL) *string {
	return buildURL("https://pub.dev/packages/%s/versions/%s", purl.Name, purl.Version)
}

// buildPypiURL constructs a PyPI package URL from a purl.
func buildPypiURL(purl packageurl.PackageURL) *string {
	return buildURL("https://pypi.org/project/%s/%s/", purl.Name, purl.Version)
}

// buildGithubURL constructs a GitHub package URL from a purl.
func buildGithubURL(purl packageurl.PackageURL) *string {
	return buildURL("https://github.com/%s/%s/tree/%s", purl.Namespace, purl.Name, purl.Version)
}

// buildDockerHubURL constructs a Docker Hub URL for docker/oci images.
// Official images (library namespace) use the "_" prefix, others use "r/" prefix.
func buildDockerHubURL(purl packageurl.PackageURL) *string {
	if purl.Namespace != "" && purl.Namespace != "library" {
		return buildURL("https://hub.docker.com/r/%s/%s", purl.Namespace, purl.Name)
	}
	return buildURL("https://hub.docker.com/_/%s", purl.Name)
}

// buildDebURL constructs a Debian package URL from a purl.
// For simplicity, we're not considering the distribution name in the URL.
func buildDebURL(purl packageurl.PackageURL) *string {
	return buildURL("https://packages.debian.org/%s", purl.Name)
}

// buildRpmURL constructs a RPM package URL from a purl.
func buildRpmURL(purl packageurl.PackageURL) *string {
	return buildURL("https://rpmfind.net/linux/rpm2html/search.php?query=%s", purl.Name)
}

// buildApkURL constructs an APK package URL from a purl.
// Search is used here because we may not know the architecture of the package.
func buildApkURL(purl packageurl.PackageURL) *string {
	return buildURL("https://pkgs.alpinelinux.org/packages?name=%s", purl.Name)
}

// buildHexURL constructs a Hex package URL from a purl.
func buildHexURL(purl packageurl.PackageURL) *string {
	return buildURL("https://hex.pm/packages/%s/%s", purl.Name, purl.Version)
}

// buildCocoapodsURL constructs a CocoaPods package URL from a purl.
func buildCocoapodsURL(purl packageurl.PackageURL) *string {
	// TODO: version is not used here, but it is possible
	return buildURL("https://cocoapods.org/pods/%s", purl.Name)
}

// buildCondaURL constructs a Conda package URL from a purl.
func buildCondaURL(purl packageurl.PackageURL) *string {
	if purl.Namespace != "" {
		return buildURL("https://anaconda.org/%s/%s", purl.Namespace, purl.Name)
	}
	return buildURL("https://anaconda.org/anaconda/%s", purl.Name)
}

// buildBitbucketURL constructs a Bitbucket package URL from a purl.
func buildBitbucketURL(purl packageurl.PackageURL) *string {
	return buildURL("https://bitbucket.org/%s/%s/src/%s", purl.Namespace, purl.Name, purl.Version)
}
