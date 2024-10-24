//go:build !windows && !unix
// +build !windows,!unix

package busybox

import (
	"errors"

	"lesiw.io/cmdio"
)

func runner() (*cmdio.Runner, error) {
	return nil, errors.New("busybox not found")
}
