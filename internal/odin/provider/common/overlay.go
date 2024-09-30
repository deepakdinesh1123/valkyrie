package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func OverlayStore(dir string, odinNixStore string) error {
	upperStore := filepath.Join(dir, "upper")
	mergedStore := filepath.Join(dir, "merged")
	workDir := filepath.Join(dir, "work")
	err := os.Mkdir(upperStore, 0755)
	if err != nil {
		return err
	}
	err = os.Mkdir(mergedStore, 0755)
	if err != nil {
		return err
	}
	err = os.Mkdir(workDir, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("mount", "-t", "overlay", "overlay",
		"-o", fmt.Sprintf("lowerdir=%s", odinNixStore),
		"-o", fmt.Sprintf("upperdir=%s", upperStore),
		"-o", fmt.Sprintf("workdir=%s", workDir),
		mergedStore,
	)

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func Cleanup(dir string) error {
	cmd := exec.Command("umount", filepath.Join(dir, "merged"))
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("rm", "-rf", dir)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
