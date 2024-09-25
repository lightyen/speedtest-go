//go:build linux && arm64

package speedtest

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	archive = "ookla-speedtest-1.2.0-linux-aarch64.tgz"
	target  = "speedtest"
)

func extract(r io.Reader) (out string, err error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}

	tr := tar.NewReader(gr)

	for h, err := tr.Next(); err == nil; h, err = tr.Next() {
		if h.Name == target {
			out = filepath.Join(os.TempDir(), target)
			if err = saveFile(tr, out); err != nil {
				return "", err
			}
			break
		}
	}

	if out == "" {
		return "", errors.New("speedtest executable not found")
	}

	return
}
