package terminal

import (
	"os/exec"

	"github.com/creack/pty"
)

func bashShell(tty *TTY) (*TTY, error) {
	cmd := exec.Command("bash")
	f, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: tty.Rows,
		Cols: tty.Columns,
	})
	if err != nil {
		tty.shell = f
	}
	return tty, err
}
