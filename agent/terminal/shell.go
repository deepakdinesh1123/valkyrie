package terminal

import (
	"fmt"
	"os/exec"

	"github.com/creack/pty"
)

func bashShell(tty *TTY) (*TTY, error) {
	cmd := exec.Command("bash")
	f, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: tty.rows,
		Cols: tty.columns,
	})
	if err != nil {
		return nil, err
	}
	tty.shell = f
	if cmd.Process == nil {
		return nil, fmt.Errorf("could not start the bash command")
	}
	tty.cmd = cmd
	return tty, err
}
