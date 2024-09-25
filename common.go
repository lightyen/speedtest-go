package speedtest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type execOptions struct {
	Ctx          context.Context
	Env          []string
	Dir          string
	InBackground bool
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
}

func execWithOptions(options execOptions, cmdstr string, args ...string) (stdoutBuf *bytes.Buffer, err error) {
	var cmd *exec.Cmd

	if options.Ctx != nil {
		cmd = exec.CommandContext(options.Ctx, cmdstr, args...)
	} else {
		cmd = exec.Command(cmdstr, args...)
	}

	cmd.Env = options.Env
	cmd.Dir = options.Dir

	var stderrBuf *bytes.Buffer

	if !options.InBackground {
		cmd.Stdin = options.Stdin
		if options.Stdout != nil {
			cmd.Stdout = options.Stdout
		} else {
			stdoutBuf = new(bytes.Buffer)
			cmd.Stdout = stdoutBuf
		}
		if options.Stderr != nil {
			cmd.Stderr = options.Stderr
		} else {
			stderrBuf = new(bytes.Buffer)
			cmd.Stderr = stderrBuf
		}
	}

	err = cmd.Start() // non-blocking
	if err != nil {
		return stdoutBuf, fmt.Errorf("[Error] EXEC %s %.64s %w", cmdstr, strings.Join(args, " "), err)
	}

	if options.InBackground {
		return stdoutBuf, nil
	}

	err = cmd.Wait()
	if err != nil || !cmd.ProcessState.Success() {
		return stdoutBuf, fmt.Errorf("[Error] EXEC %s %.64s %.64s", cmdstr, strings.Join(args, " "), stderrBuf.String())
	}

	return stdoutBuf, nil
}

// https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-macosx-universal.tgz
// https://install.speedtest.net/app/cli/ookla-speedtest-1.2.0-linux-armhf.tgz

const api = "https://install.speedtest.net/app/cli/"

var lc sync.Mutex

func ensure() (name string, err error) {
	lc.Lock()
	defer lc.Unlock()

	name = filepath.Join(os.TempDir(), target)
	if info, err := os.Stat(name); err == nil {
		if info.Mode().IsRegular() {
			return name, nil
		}
	}

	resp, err := http.Get(api + archive)
	if err != nil {
		return "", err
	}

	defer func() {
		if err2 := resp.Body.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()

	return extract(resp.Body)
}

func saveFile(r io.Reader, name string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}

	defer func() {
		if err2 := f.Close(); err == nil && err2 != nil {
			err = err2
		}
	}()
	_, err = io.Copy(f, r)
	return err
}
