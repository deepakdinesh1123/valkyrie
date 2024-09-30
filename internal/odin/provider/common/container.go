package common

import (
	"os/exec"
	"strconv"
)

func KillContainer(pid int) error {
	cmd := exec.Command("kill", "-KILL", strconv.Itoa(pid))
	return cmd.Run()
}
