# lesiw.io/cmdio/x/busybox

Package busybox provides a `cmdio.Runner` that runs commands in a `busybox`.

The runner prefers existing installs of `busybox` if present, otherwise it will
attempt to download a static build for the current `runtime.GOOS` and
`runtime.GOARCH`.

If no install of `busybox` is available and the underlying system is a
unix-like, the runner will "fall through" to a normal `sys.Runner()`.
