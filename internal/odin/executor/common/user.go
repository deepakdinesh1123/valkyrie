package common

import (
	"fmt"
	"os"
	"os/user"
)

func GetUserInfo() (*user.User, error) {
	eUser := os.Getenv("SUDO_USER")
	userInfo, err := user.Lookup(eUser)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup user info: %w", err)
	}
	return userInfo, nil
}
