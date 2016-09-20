package csscleaner

import (
	"lib"
)

// Cleanup cleans up the given css via csstidy
// csstidy must be installed and inside the path
func Cleanup(css string) (string, error) {
	return lib.Cleanup(css)
}
