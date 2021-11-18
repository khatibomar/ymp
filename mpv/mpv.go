package mpv

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	tmpdir = filepath.Join(os.TempDir(), "ymp")
)

type Mpv struct {
	Cmd        *exec.Cmd
	SocketPath string
}

func New() (*Mpv, error, func()) {
	sockPath := filepath.Join(tmpdir, "mpv", "mpv.sock")
	if err := os.MkdirAll(filepath.Dir(sockPath), os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to make socket directory"), func() {}
	}

	args := []string{
		"--idle",
		"--quiet",
		"--pause",
		"--no-input-terminal",
		"--loop-playlist=no",
		"--gapless-audio=weak",
		"--replaygain=no",
		// "--replaygain-clip=no",
		"--ad=lavc:*",
		"--input-ipc-server=" + sockPath,
		"--volume=100",
		"--volume-max=100",
		"--no-video",
	}

	cmd := exec.Command("mpv", args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr

	cleanup := func() {
		cmd.Process.Kill()
		os.RemoveAll(sockPath)
	}
	return &Mpv{
		Cmd:        cmd,
		SocketPath: sockPath,
	}, nil, cleanup
}
