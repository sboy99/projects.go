package version

import "fmt"

// Build-time variables (set via ldflags)
var (
	BuildTime    string
	BuildVersion string
	BuildCommit  string
)

// GetVersion returns the version information
func GetVersion() string {
	return BuildVersion
}

// GetBuildInfo returns build information
func GetBuildInfo() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Time: %s", BuildVersion, BuildCommit, BuildTime)
}

// String returns a formatted version string
func String() string {
	return GetBuildInfo()
}
