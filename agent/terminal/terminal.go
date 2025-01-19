package terminal

import (
	"os"

	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
)

// TTY represents a terminal with customizable TTYOptss.
type TTY struct {
	Rows    uint16
	Columns uint16

	shell *os.File
}

// TTYOpts defines a function type for applying options to a TTY.
type TTYOpts func(*TTY)

// NewTTY initializes a TTY with the provided options
func NewTTY(terminal *schemas.Newterminal, opts ...TTYOpts) (*TTY, error) {
	// Default values for the TTY.
	t := &TTY{
		Rows:    24,
		Columns: 80,
	}

	// Apply each TTYOpts to the TTY.
	for _, opt := range opts {
		opt(t)
	}

	switch terminal.Shell {
	case "bash":
		{
			t, err := bashShell(t)
			if err != nil {
				return t, err
			}
		}
	}

	return t, nil
}

func WithRows(rows uint16) TTYOpts {
	return func(t *TTY) {
		t.Rows = rows
	}
}

func WithCols(cols uint16) TTYOpts {
	return func(t *TTY) {
		t.Columns = cols
	}
}
