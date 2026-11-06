// Package version holds the single source of truth for the Hemlis version.
package version

// Version is the release version. Bump it here when cutting a release.
const Version = "0.2.0"

// Commit and BuildDate are injected at build time via -ldflags (see build.sh).
var (
	Commit    = "unknown"
	BuildDate = "unknown"
)
