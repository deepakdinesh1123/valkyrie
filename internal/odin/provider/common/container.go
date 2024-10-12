package common

import (
	"fmt"
	"os/exec"
	"strconv"
)

func KillContainer(pid int) error {
	if pid == 0 {
		return fmt.Errorf("pid is 0")
	}
	fmt.Printf("Killing container with pid %d\n", pid)
	cmd := exec.Command("kill", "-KILL", strconv.Itoa(pid))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill container: %w", err)
	}
	return nil
}
