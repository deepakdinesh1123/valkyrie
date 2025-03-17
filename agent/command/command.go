package command

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
	"github.com/google/uuid"
)

type Command struct {
	rows    uint16
	columns uint16

	shell *os.File
	cmd   *exec.Cmd

	Completed bool
}

type CommandOpts func(*Command)

func WithRows(rows uint16) CommandOpts {
	return func(t *Command) {
		t.rows = rows
	}
}

func WithCols(cols uint16) CommandOpts {
	return func(t *Command) {
		t.columns = cols
	}
}

func NewCommand(ec *schemas.ExecuteCommand, opts ...CommandOpts) (*Command, string, error) {
	cmd := &Command{
		rows:    24,
		columns: 80,
	}

	// Apply each CommandOpts to the Command.
	for _, opt := range opts {
		opt(cmd)
	}

	command := strings.Fields(ec.Command)
	var cmdExec *exec.Cmd
	if len(command) > 1 {
		cmdExec = exec.Command(command[0])
	} else {
		cmdExec = exec.Command(command[0], command[:1]...)
	}

	if ec.WorkDir != nil {
		if _, err := os.Stat(*ec.WorkDir); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(*ec.WorkDir, 0755)
				if err != nil {
					return nil, "", err
				}
			}
		}
		cmdExec.Dir = *ec.WorkDir
	}

	f, err := pty.StartWithSize(cmdExec, &pty.Winsize{
		Rows: cmd.rows,
		Cols: cmd.columns,
	})

	if err != nil {
		return nil, "", err
	}

	cmd.shell = f
	cmd.cmd = cmdExec

	cmdId := uuid.NewString()
	return cmd, cmdId, nil
}

// Read reads from the terminal output.
func (cmd *Command) Read() (string, error) {
	buf := make([]byte, 1024)
	n, err := cmd.shell.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

// Write sends input to the terminal.
func (cmd *Command) Write(input []byte) error {
	_, err := cmd.shell.Write(input)
	return err
}

// Terminates the Command session.
func (cmd *Command) Terminate() error {
	if cmd.shell != nil {
		cmd.shell.Close()
	}
	if cmd.cmd != nil && cmd.cmd.Process != nil {
		return cmd.cmd.Process.Signal(syscall.SIGKILL)
	}
	return nil
}

func (cmd *Command) Wait() error {
	err := cmd.cmd.Wait()
	return err
}
