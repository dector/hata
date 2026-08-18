package version

import (
	"os"
	"strings"
)

// Version is the application version.
//
// It defaults to "dev" and can be overridden at build time with:
//
//	go build -ldflags "-X hata/version.Version=<version>" ./cmd/hata
var Version = "dev"

// Get returns the effective application version.
//
// HATA_VERSION can override the build-time version when running the app.
func Get() string {
	if value := strings.TrimSpace(os.Getenv("HATA_VERSION")); value != "" {
		return value
	}

	if value := strings.TrimSpace(Version); value != "" {
		return value
	}

	return "dev"
}
