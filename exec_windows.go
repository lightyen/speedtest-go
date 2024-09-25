//go:build amd64

package speedtest

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	archive = "ookla-speedtest-1.2.0-win64.zip"
	target  = "speedtest.exe"
)

func extract(reader io.Reader) (out string, err error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	f, err := r.Open(target)
	if err != nil {
		return "", err
	}
	defer func() {
		if err2 := f.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()

	out = filepath.Join(os.TempDir(), target)
	if err = saveFile(f, out); err != nil {
		return "", err
	}

	if out == "" {
		return "", errors.New("speedtest executable not found")
	}

	return
}
