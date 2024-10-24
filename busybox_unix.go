//go:build unix
// +build unix

package busybox

import (
	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

func runner() (*cmdio.Runner, error) {
	// Fall through to the system implementation for now.
	// TODO: Download static busybox builds where possible.
	return sys.Runner(), nil
}
