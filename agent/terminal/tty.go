package terminal

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
	"github.com/google/uuid"
)

// TTY represents a terminal with customizable TTYOpts.
type TTY struct {
	rows    uint16
	columns uint16

	shell *os.File
	cmd   *exec.Cmd
}

// TTYOpts defines a function type for applying options to a TTY.
type TTYOpts func(*TTY)

// NewTTY initializes a TTY with the provided options
func NewTTY(terminal *schemas.NewTerminal, opts ...TTYOpts) (*TTY, string, error) {
	// Default values for the TTY.
	t := &TTY{
		rows:    24,
		columns: 80,
	}

	// Apply each TTYOpts to the TTY.
	for _, opt := range opts {
		opt(t)
	}

	t, err := bashShell(t)
	if err != nil {
		return t, "", err
	}

	tid := uuid.NewString()

	return t, tid, nil
}

func WithRows(rows uint16) TTYOpts {
	return func(t *TTY) {
		t.rows = rows
	}
}

func WithCols(cols uint16) TTYOpts {
	return func(t *TTY) {
		t.columns = cols
	}
}

// Read reads from the terminal output.
func (tty *TTY) Read() (string, error) {
	buf := make([]byte, 1024)
	n, err := tty.shell.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

// Write sends input to the terminal.
func (tty *TTY) Write(input []byte) error {
	_, err := tty.shell.Write(input)
	return err
}

// Close terminates the TTY session.
func (tty *TTY) Close() error {
	if tty.shell != nil {
		tty.shell.Close()
	}
	if tty.cmd != nil && tty.cmd.Process != nil {
		return tty.cmd.Process.Signal(syscall.SIGKILL)
	}
	return nil
}
