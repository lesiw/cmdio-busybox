# lesiw.io/cmdio/x/busybox

Package busybox provides a `cmdio.Runner` that acts as an additional
best-effort compatibility layer for `cmdio` automation.

The runner prefers existing installs of `busybox` if present, otherwise it will
attempt to download a static build for the current `runtime.GOOS` and
`runtime.GOARCH`.

If no install of `busybox` is available and the underlying system is a
unix-like, commands will fall through to the native system.
