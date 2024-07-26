package cmd_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/cmd"
)

func TestRootCommand(t *testing.T) {
	b := bytes.NewBufferString("")
	cmd.RootCmd.SetOut(b)
	cmd.RootCmd.Execute()
	out, err := io.ReadAll(b)
	if err != nil {
		t.Errorf("Failed to read output: %v", err)
	}
	if string(out) != "Execute odin --help for more information on using odin.\n" {
		t.Errorf("Unexpected output: %s", string(out))
	}
}
