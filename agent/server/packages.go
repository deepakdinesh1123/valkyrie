package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
)

func (s *Server) handleInstallNixPackage(ctx context.Context, c *websocket.Conn, data []byte) {
	var inp schemas.InstallNixPackage

	msgType := "InstallNixPackageResponse"

	sendResponse := func(success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.InstallNixPackageResponse{
			MsgType: &msgType,
			Success: success,
			Msg:     msg,
		})
	}

	// Decode the incoming data into the inp variable
	err := json.Unmarshal(data, &inp)
	if err != nil {
		sendResponse(false, fmt.Sprintf("Failed to decode input data: %v", err))
		return
	}

	// Create a ticker to send messages every 5 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Create a channel to signal when the installation is done
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				SendJSONMessage(ctx, c, schemas.InstallNixPackageResponse{
					MsgType: &msgType,
					Success: false,
					Msg:     "Package is being installed...",
				})
			case <-done:
				return
			}
		}
	}()

	cmd := exec.Command("nix", "profile", "install", fmt.Sprintf("%s#%s", *inp.Channel, inp.PkgName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		sendResponse(false, fmt.Sprintf("Failed to install package: %v. Output: %s", err, string(output)))
		close(done)
		return
	}

	// Check the exit code of the command
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
		sendResponse(false, fmt.Sprintf("Failed to install package: exit code %d. Output: %s", exitErr.ExitCode(), string(output)))
		close(done)
		return
	}

	sendResponse(true, "Package installed successfully")
	close(done)
}

func (s *Server) handleUninstallNixPackage(ctx context.Context, c *websocket.Conn, data []byte) {
	var inp schemas.UninstallNixPackage

	msgType := "UninstallNixPackageResponse"

	sendResponse := func(success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.UninstallNixPackageResponse{
			MsgType: &msgType,
			Success: success,
			Msg:     msg,
		})
	}

	// Decode the incoming data into the inp variable
	err := json.Unmarshal(data, &inp)
	if err != nil {
		sendResponse(false, fmt.Sprintf("Failed to decode input data: %v", err))
		return
	}

	cmd := exec.Command("nix", "profile", "remove", inp.PkgName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		sendResponse(false, fmt.Sprintf("Failed to uninstall package: %v. Output: %s", err, string(output)))
		return
	}

	// Check the exit code of the command
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
		sendResponse(false, fmt.Sprintf("Failed to uninstall package: exit code %d. Output: %s", exitErr.ExitCode(), string(output)))
		return
	}

	sendResponse(true, "Package uninstalled successfully")
}
