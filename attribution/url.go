package attribution

import (
	"fmt"

	"github.com/package-url/packageurl-go"
)

// PurlToURL constructs a package management URL from a purl string.
// Returns nil if the purl cannot be parsed or the type is not supported.
func PurlToURL(purlString string) *string {
	if purlString == "" {
		return nil
	}

	purl, err := packageurl.FromString(purlString)
	if err != nil {
		logger.Warn("failed to parse purl", "purl", purlString, "error", err) //nolint:sloglint // Package-level logger
		return nil
	}

	return mapPurlToURL(purl)
}

func mapPurlToURL(purl packageurl.PackageURL) *string {
	// See https://github.com/package-url/purl-spec#known-purl-types
	switch purl.Type {
	case "cargo":
		return buildURL("https://crates.io/crates/%s/%s", purl.Name, purl.Version)
	case "composer":
		return buildURL("https://packagist.org/packages/%s/%s#%s", purl.Namespace, purl.Name, purl.Version)
	case "gem":
		return buildURL("https://rubygems.org/gems/%s/versions/%s", purl.Name, purl.Version)
	case "golang":
		if purl.Namespace != "" {
			return buildURL("https://pkg.go.dev/%s/%s@%s", purl.Namespace, purl.Name, purl.Version)
		}
		return buildURL("https://pkg.go.dev/%s@%s", purl.Name, purl.Version)
	case "maven":
		return buildURL("https://mvnrepository.com/artifact/%s/%s/%s", purl.Namespace, purl.Name, purl.Version)
	case "npm":
		if purl.Namespace != "" {
			return buildURL("https://www.npmjs.com/package/%s/%s/v/%s", purl.Namespace, purl.Name, purl.Version)
		}
		return buildURL("https://www.npmjs.com/package/%s/v/%s", purl.Name, purl.Version)
	case "nuget":
		return buildURL("https://www.nuget.org/packages/%s/%s", purl.Name, purl.Version)
	case "pub":
		return buildURL("https://pub.dev/packages/%s/versions/%s", purl.Name, purl.Version)
	case "pypi":
		return buildURL("https://pypi.org/project/%s/%s/", purl.Name, purl.Version)
	case "github":
		return buildURL("https://github.com/%s/%s/tree/%s", purl.Namespace, purl.Name, purl.Version)
	case "docker", "oci":
		return buildDockerHubURL(purl.Namespace, purl.Name)
	case "deb":
		return buildURL("https://packages.debian.org/%s", purl.Name)
	case "rpm":
		return buildURL("https://rpmfind.net/linux/rpm2html/search.php?query=%s", purl.Name)
	case "apk":
		return buildURL("https://pkgs.alpinelinux.org/packages?name=%s", purl.Name)
	case "hex":
		return buildURL("https://hex.pm/packages/%s/%s", purl.Name, purl.Version)
	case "cocoapods":
		return buildURL("https://cocoapods.org/pods/%s", purl.Name)
	case "conda":
		if purl.Namespace != "" {
			return buildURL("https://anaconda.org/%s/%s", purl.Namespace, purl.Name)
		}
		return buildURL("https://anaconda.org/anaconda/%s", purl.Name)
	case "bitbucket":
		return buildURL("https://bitbucket.org/%s/%s/src/%s", purl.Namespace, purl.Name, purl.Version)
	case "alpm", "bitnami", "conan", "cran",
		"generic", "hackage", "huggingface", "mlflow",
		"qpkg", "swid", "swift":
		logger.Debug("purl type not yet supported", "type", purl.Type) //nolint:sloglint // Package-level logger
		return nil
	default:
		logger.Warn("unknown purl type", "type", purl.Type) //nolint:sloglint // Package-level logger
		return nil
	}
}

// buildURL constructs a URL from a format string and arguments.
func buildURL(format string, args ...interface{}) *string {
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
