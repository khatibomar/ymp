package mpv

import (
	"bufio"
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
	OutBuff    *bufio.Reader
}

func New() (*Mpv, error, func()) {
	sockPath := filepath.Join(tmpdir, "mpv", "mpv.sock")
	if err := os.MkdirAll(filepath.Dir(sockPath), os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to make socket directory"), func() {}
	}

	args := []string{
		"--idle",
		"--pause",
		"--no-input-terminal",
		"--loop-playlist=no",
		"--gapless-audio=weak",
		"--replaygain=no",
		// "--replaygain-clip=no",
		"--input-ipc-server=" + sockPath,
		"--volume=100",
		"--volume-max=100",
		"--no-video",
		"--no-resume-playback",
	}

	cmd := exec.Command("mpv", args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr

	stdout, _ := cmd.StdoutPipe()

	buf := bufio.NewReader(stdout)

	return &Mpv{
			Cmd:        cmd,
			SocketPath: sockPath,
			OutBuff:    buf,
		}, nil, func() {
			cmd.Process.Kill()
			os.RemoveAll(sockPath)
		}
}
