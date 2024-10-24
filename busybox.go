package busybox

import (
	"fmt"
	"io"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sub"
	"lesiw.io/cmdio/sys"
)

func Runner() (*cmdio.Runner, error) {
	rnr := sys.Runner()
	if _, err := io.ReadAll(rnr.Command("busybox")); err == nil {
		return sub.New("busybox"), nil
	}
	if rnr, err := runner(); err != nil {
		return nil, fmt.Errorf("failed to fetch busybox: %w", err)
	} else {
		return rnr, nil
	}
}
