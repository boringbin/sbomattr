# sbomattr

A simple CLI tool to create an aggregated notice for one or more SBOMs
([SPDX](https://spdx.dev/) or [CycloneDX](https://cyclonedx.org/)).

## Usage

```text
Usage: sbomattr [OPTIONS] <file-or-directory>...

Create an aggregated notice for one or more SBOMs.

Arguments:
  file-or-directory   SBOM files or directories containing SBOM files

Options:
  -v    Verbose output (debug mode)
  -version
        Show version and exit
```

## Why?

Provide clear attribution for software dependencies in a simple, verifiable format.

When distributing software (especially closed source), you could want to aggregate license information from multiple
SBOMs into a single notice file. This tool does one thing well: combine SBOMs into unified attribution notices.

## What is the `URL` Field?

The `URL` field is the quickest way to validate the package information for people who don't care about
[the purl specification](https://github.com/package-url/purl-spec).

If your SBOM already contains a relevant URL field, this will be used. Otherwise, the `Purl` field will be used to
construct a URL.

## Supported Formats

- [SPDX 2.3](https://spdx.github.io/spdx-spec/v2.3/) (JSON)
- [CycloneDX 1.4](https://cyclonedx.org/docs/1.4/json/) (JSON)
- GitHub-wrapped SBOMs (JSON)

## License

[MIT](LICENSE)
