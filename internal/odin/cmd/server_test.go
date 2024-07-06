package cmd_test

import (
	"bytes"
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/cmd"
)

func TestServerCommand(t *testing.T) {
	b := bytes.NewBufferString("")
	cmd.ServerCmd.SetOut(b)
	cmd.ServerCmd.Execute()
}
