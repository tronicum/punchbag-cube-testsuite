package sim

import "os"

// IsDebug returns true if debug mode is enabled via env var
func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}
