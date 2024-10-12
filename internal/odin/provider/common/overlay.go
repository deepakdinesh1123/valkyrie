package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func OverlayStore(dir string, odinNixStore string) error {
	userInfo, err := GetUserInfo()
	if err != nil {
		return fmt.Errorf("failed to get user name: %w", err)
	}

	groupId, err := strconv.Atoi(userInfo.Gid)
	if err != nil {
		return fmt.Errorf("failed to convert group id to int: %w", err)
	}
	userId, err := strconv.Atoi(userInfo.Uid)
	if err != nil {
		return fmt.Errorf("failed to convert user id to int: %w", err)
	}

	upperStore := filepath.Join(dir, "upper")
	mergedStore := filepath.Join(dir, "merged")
	workDir := filepath.Join(dir, "work")
	err = os.Mkdir(upperStore, 0755)
	if err != nil {
		return fmt.Errorf("failed to create upper store: %w", err)
	}
	err = os.Mkdir(mergedStore, 0755)
	if err != nil {
		return fmt.Errorf("failed to create merged store: %w", err)
	}
	err = os.Mkdir(workDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create work dir: %w", err)
	}

	// change the ownership of the store and work dir to the user
	err = os.Chown(upperStore, userId, groupId)
	if err != nil {
		return fmt.Errorf("failed to change ownership of upper store: %w", err)
	}
	err = os.Chown(mergedStore, userId, groupId)
	if err != nil {
		return fmt.Errorf("failed to change ownership of merged store: %w", err)
	}
	err = os.Chown(workDir, userId, groupId)
	if err != nil {
		return fmt.Errorf("failed to change ownership of work dir: %w", err)
	}

	// mount the overlay
	cmd := exec.Command("mount", "-t", "overlay", "overlay",
		"-o", fmt.Sprintf("lowerdir=%s", odinNixStore),
		"-o", fmt.Sprintf("upperdir=%s", upperStore),
		"-o", fmt.Sprintf("workdir=%s", workDir),
		mergedStore,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount overlay: %w: %s", err, string(out))
	}

	// change the ownership of the merged store to the user
	err = os.Chown(mergedStore, userId, groupId)
	if err != nil {
		return fmt.Errorf("failed to change ownership of merged store: %w", err)
	}
	err = os.Chmod(mergedStore, 0777)
	if err != nil {
		return fmt.Errorf("failed to change permissions of merged store: %w", err)
	}
	return nil
}

func Cleanup(dir string) error {
	cmd := exec.Command("umount", filepath.Join(dir, "merged"))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unmount overlay: %w", err)
	}
	cmd = exec.Command("rm", "-rf", dir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}
	return nil
}
