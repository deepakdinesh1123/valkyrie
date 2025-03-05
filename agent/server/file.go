package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
)

func (s *Server) handleUpsertFile(ctx context.Context, c *websocket.Conn, data []byte) {
	var uf schemas.UpsertFile
	msgType := "UpsertFileResponse"

	// Unmarshal the incoming data
	if err := json.Unmarshal(data, &uf); err != nil {
		SendJSONMessage(ctx, c, schemas.UpsertFileResponse{
			MsgType: &msgType,
			Success: false,
			Msg:     fmt.Sprintf("error unmarshaling request: %v", err),
			Path:    uf.Path,
		})
		return
	}

	sendResponse := func(success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.UpsertFileResponse{
			MsgType: &msgType,
			Success: success,
			Msg:     msg,
			Path:    uf.Path,
		})
	}

	// Create directory structure if it doesn't exist
	dir := filepath.Dir(uf.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		sendResponse(false, fmt.Sprintf("error creating directory structure: %v", err))
		return
	}

	// Check if the file exists
	if _, err := os.Stat(uf.Path); err == nil {
		if !os.IsNotExist(err) {
			// File doesn't exist, create it
			_, err = os.Create(uf.Path)
			if err != nil {
				sendResponse(false, fmt.Sprintf("error creating file: %v", err))
				return
			}
		} else {
			// Other error occurred
			sendResponse(false, fmt.Sprintf("error checking file status: %v", err))
			return
		}
	}

	// Handle file content based on whether a patch is provided
	if uf.Patch != "" {
		// Create temporary files for the patch
		patchFile, err := os.CreateTemp("", "patch-*.diff")
		if err != nil {
			sendResponse(false, fmt.Sprintf("error creating temporary patch file: %v", err))
			return
		}
		defer os.Remove(patchFile.Name())

		// Write the patch to the temporary file
		if _, err := patchFile.Write([]byte(uf.Patch)); err != nil {
			sendResponse(false, fmt.Sprintf("error writing to patch file: %v", err))
			return
		}
		patchFile.Close() // Close to ensure content is flushed

		// Execute the patch command
		cmd := exec.Command("patch", uf.Path, patchFile.Name())
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			sendResponse(false, fmt.Sprintf("error applying patch: %v, stderr: %s", err, stderr.String()))
			return
		}
	} else if uf.Content != "" {
		// Replace the entire file content
		file, err := os.Create(uf.Path) // Create/truncate file
		if err != nil {
			sendResponse(false, fmt.Sprintf("error creating/truncating file: %v", err))
			return
		}
		defer file.Close()

		// Write new content
		if _, err := file.Write([]byte(uf.Content)); err != nil {
			sendResponse(false, fmt.Sprintf("error writing content: %v", err))
			return
		}
	} else {
		sendResponse(true, "File created")
		return
	}

	sendResponse(true, "file updated successfully")
}

func (s *Server) handleReadFile(ctx context.Context, c *websocket.Conn, data []byte) {
	var rf schemas.ReadFile
	msgType := "ReadFileResponse"

	// Unmarshal the incoming data
	if err := json.Unmarshal(data, &rf); err != nil {
		SendJSONMessage(ctx, c, schemas.ReadFileResponse{
			MsgType: &msgType,
			Success: false,
			Msg:     fmt.Sprintf("error unmarshaling request: %v", err),
			Path:    rf.Path,
		})
		return
	}

	sendResponse := func(success bool, msg string, content string) {
		SendJSONMessage(ctx, c, schemas.ReadFileResponse{
			MsgType: &msgType,
			Success: success,
			Msg:     msg,
			Path:    rf.Path,
			Content: content,
		})
	}

	// Check if the file exists
	if _, err := os.Stat(rf.Path); os.IsNotExist(err) {
		sendResponse(false, "file does not exist", "")
		return
	} else if err != nil {
		sendResponse(false, fmt.Sprintf("error checking file status: %v", err), "")
		return
	}

	// Read the file content
	content, err := os.ReadFile(rf.Path)
	if err != nil {
		sendResponse(false, fmt.Sprintf("error reading file: %v", err), "")
		return
	}

	sendResponse(true, "file read successfully", string(content))
}

func (s *Server) handleDeleteFile(ctx context.Context, c *websocket.Conn, data []byte) {
	var df schemas.DeleteFile
	msgType := "DeleteFileResponse"

	// Unmarshal the incoming data
	if err := json.Unmarshal(data, &df); err != nil {
		SendJSONMessage(ctx, c, schemas.DeleteFileResponse{
			MsgType: &msgType,
			Success: false,
			Msg:     fmt.Sprintf("error unmarshaling request: %v", err),
			Path:    df.Path,
		})
		return
	}

	sendResponse := func(success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.DeleteFileResponse{
			MsgType: &msgType,
			Success: success,
			Msg:     msg,
			Path:    df.Path,
		})
	}

	// Check if the file exists
	if _, err := os.Stat(df.Path); os.IsNotExist(err) {
		sendResponse(false, "file does not exist")
		return
	} else if err != nil {
		sendResponse(false, fmt.Sprintf("error checking file status: %v", err))
		return
	}

	// Delete the file
	if err := os.Remove(df.Path); err != nil {
		sendResponse(false, fmt.Sprintf("error deleting file: %v", err))
		return
	}

	sendResponse(true, "file deleted successfully")
}
