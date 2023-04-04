// The version package provides a location to set the release versions for all
// packages to consume, without creating import cycles.
//
// This package should not import any other terraform packages.
package version

import (
	_ "embed"
	"os"
	"strings"

	version "github.com/hashicorp/go-version"
)

// The main version number that is being run at the moment.
var Version string

// A pre-release marker for the version. If this is "" (empty string)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
var Prerelease string

// SemVer is an instance of version.Version. This has the secondary
// benefit of verifying during tests and init time that our version is a
// proper semantic version, which should always be the case.
var SemVer *version.Version

// rawVersion is the current version as a string, as read from the VERSION
// file. This must be a valid semantic version.
//
//go:embed VERSION
var rawVersion string

func init() {
	rawVersion = strings.TrimSpace(rawVersion)
	version.Must(version.NewVersion(rawVersion))
	Version, Prerelease = extractPrerelease(rawVersion)
	SemVer = version.Must(version.NewVersion(Version))
}

// Header is the header name used to send the current terraform version
// in http requests.
const Header = "Terraform-Version"

// String returns the complete version string, including prerelease
func String() string {
	return SemVer.String()
}

func extractPrerelease(rawVersion string) (string, string) {
	parts := strings.Split(rawVersion, "-")

	var prerelease string
	if os.Getenv("TFDEV") == "1" {
		prerelease = "dev"
	} else if len(parts) > 1 {
		prerelease = parts[1]
	} else {
		prerelease = ""
	}

	return parts[0], prerelease
}
