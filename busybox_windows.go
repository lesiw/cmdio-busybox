//go:build windows
// +build windows

package busybox

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/crypto/openpgp"
	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sub"
)

const (
	w32exe = "https://frippery.org/files/busybox/"
	w32sum = "https://frippery.org/files/busybox/SHA256SUM"
	w32sig = "https://frippery.org/files/busybox/SHA256SUM.sig"
)

//go:embed frippery.asc
var w32key []byte

func runner() (*cmdio.Runner, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user cache directory: %w", err)
	}
	cache = filepath.Join(cache, "cmdio")
	if err = os.MkdirAll(cache, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory %q: %w",
			cache, err)
	}
	path := filepath.Join(cache, "busybox.exe")
	f, err := os.Stat(path)
	if err == nil && !f.IsDir() {
		return sub.New(path), nil
	}
	hashes, err := fetchHashes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch busybox-w32 hashes: %w", err)
	}
	file := "busybox.exe"
	switch runtime.GOARCH {
	case "arm64":
		file = "busybox64a.exe"
	case "amd64":
		file = "busybox64.exe"
	}
	if err := dl(path, w32exe+file, hashes[file]); err != nil {
		return nil, err
	}
	return sub.New(path), nil
}

func dl(dst, url string, sum []byte) error {
	var tmp *os.File
	defer func() {
		if tmp == nil {
			return
		}
		tmp.Close()
		os.RemoveAll(tmp.Name())
	}()
	r, err := fetchBody(url)
	if err != nil {
		return fmt.Errorf("failed to get %q: %w", url, err)
	}
	if tmp, err = os.CreateTemp("", "cmdio_busybox_*.tmp"); err != nil {
		return fmt.Errorf("failed to create tmpfile: %w", err)
	}
	tr := io.TeeReader(r, tmp)
	hsh := sha256.New()
	if _, err := io.Copy(hsh, tr); err != nil {
		return fmt.Errorf("failed reading %q: %w", url, err)
	}
	if got, want := hsh.Sum(nil), sum; !bytes.Equal(got, want) {
		return fmt.Errorf("bad checksum for %q: %q, want %q", url, got, want)
	}
	src := tmp.Name()
	tmp.Close()
	tmp = nil
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("failed to move %q to %q: %w", src, dst, err)
	}
	return nil
}

func fetchHashes() (map[string][]byte, error) {
	hshbuf, err := fetchUrl(w32sum)
	if err != nil {
		return nil, err
	}
	sigbuf, err := fetchUrl(w32sig)
	if err != nil {
		return nil, err
	}
	key, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(w32key))
	if err != nil {
		return nil, fmt.Errorf("bad key: %w", err)
	}
	_, err = openpgp.CheckDetachedSignature(key,
		bytes.NewReader(hshbuf), bytes.NewReader(sigbuf))
	if err != nil {
		return nil, fmt.Errorf("bad signature: %w", err)
	}
	scr := bufio.NewScanner(bytes.NewReader(hshbuf))
	ret := make(map[string][]byte)
	for scr.Scan() {
		line := scr.Text()
		hash, file, ok := strings.Cut(line, "  ")
		if !ok {
			return nil, fmt.Errorf("bad sum %q", line)
		}
		buf, err := hex.DecodeString(hash)
		if err != nil {
			return nil, fmt.Errorf("bad hash %q: %w", hash, err)
		}
		ret[file] = buf
	}
	return ret, nil
}

func fetchUrl(url string) ([]byte, error) {
	r, err := fetchBody(url)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	ret, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", url, err)
	}

	return ret, nil
}

func fetchBody(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %q: %w", url, err)
	}
	req.Header.Set("User-Agent", "")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get %q: %w", url, err)
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get %q: %d", url, rsp.StatusCode)
	}

	return rsp.Body, nil
}
