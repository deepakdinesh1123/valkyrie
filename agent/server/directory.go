package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
)

func (s *Server) handleDeleteDirectory(ctx context.Context, c *websocket.Conn, data []byte) {
	var dd schemas.DeleteDirectory
	msgType := "DeleteDirectoryResponse"

	// Unmarshal the incoming data
	if err := json.Unmarshal(data, &dd); err != nil {
		SendJSONMessage(ctx, c, schemas.DeleteDirectoryResponse{
			MsgType: &msgType,
			Success: false,
			Msg:     fmt.Sprintf("error unmarshaling request: %v", err),
			Path:    dd.Path,
		})
		return
	}

	sendResponse := func(success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.DeleteDirectoryResponse{
			MsgType: &msgType,
			Success: success,
			Msg:     msg,
			Path:    dd.Path,
		})
	}

	// Check if the directory exists
	info, err := os.Stat(dd.Path)
	if os.IsNotExist(err) {
		sendResponse(false, "directory does not exist")
		return
	} else if err != nil {
		sendResponse(false, fmt.Sprintf("error checking directory status: %v", err))
		return
	}

	// Make sure it's actually a directory
	if !info.IsDir() {
		sendResponse(false, "path exists but is not a directory")
		return
	}

	// Delete the directory (recursively)
	if err := os.RemoveAll(dd.Path); err != nil {
		sendResponse(false, fmt.Sprintf("error deleting directory: %v", err))
		return
	}

	sendResponse(true, "directory deleted successfully")
}

func (s *Server) handleReadDirectory(ctx context.Context, c *websocket.Conn, data []byte) {
	var rd schemas.ReadDirectory
	msgType := "ReadDirectoryResponse"

	// Unmarshal the incoming data
	if err := json.Unmarshal(data, &rd); err != nil {
		SendJSONMessage(ctx, c, schemas.ReadDirectoryResponse{
			MsgType: &msgType,
			Success: false,
			Msg:     fmt.Sprintf("error unmarshaling request: %v", err),
			Path:    rd.Path,
		})
		return
	}

	sendResponse := func(success bool, msg string, contents string) {
		SendJSONMessage(ctx, c, schemas.ReadDirectoryResponse{
			MsgType:  &msgType,
			Success:  success,
			Msg:      msg,
			Path:     rd.Path,
			Contents: contents,
		})
	}

	// Check if the directory exists
	info, err := os.Stat(rd.Path)
	if os.IsNotExist(err) {
		sendResponse(false, "directory does not exist", "")
		return
	} else if err != nil {
		sendResponse(false, fmt.Sprintf("error checking directory status: %v", err), "")
		return
	}

	// Make sure it's actually a directory
	if !info.IsDir() {
		sendResponse(false, "path exists but is not a directory", "")
		return
	}

	// Read the directory entries
	dirEntries, err := os.ReadDir(rd.Path)
	if err != nil {
		sendResponse(false, fmt.Sprintf("error reading directory: %v", err), "")
		return
	}

	// Format directory contents as a string
	var contents strings.Builder
	for _, entry := range dirEntries {
		entryType := "f"
		if entry.IsDir() {
			entryType = "d"
		}
		contents.WriteString(fmt.Sprintf("%s %s\n", entryType, entry.Name()))
	}

	sendResponse(true, "directory read successfully", contents.String())
}

func (s *Server) handleUpsertDirectory(ctx context.Context, c *websocket.Conn, data []byte) {
	var ud schemas.UpsertDirectory
	msgType := "UpsertDirectoryResponse"

	// Unmarshal the incoming data
	if err := json.Unmarshal(data, &ud); err != nil {
		SendJSONMessage(ctx, c, schemas.UpsertDirectoryResponse{
			MsgType: &msgType,
			Success: false,
			Msg:     fmt.Sprintf("error unmarshaling request: %v", err),
			Path:    ud.Path,
		})
		return
	}

	sendResponse := func(success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.UpsertDirectoryResponse{
			MsgType: &msgType,
			Success: success,
			Msg:     msg,
			Path:    ud.Path,
		})
	}

	// Check if the directory already exists
	info, err := os.Stat(ud.Path)
	if err == nil {
		// Path exists, check if it's a directory
		if !info.IsDir() {
			sendResponse(false, "path exists but is not a directory")
			return
		}

		// Directory already exists
		sendResponse(true, "directory already exists")
		return
	} else if !os.IsNotExist(err) {
		// Some other error occurred
		sendResponse(false, fmt.Sprintf("error checking directory status: %v", err))
		return
	}

	// Create the directory
	if err := os.MkdirAll(ud.Path, 0755); err != nil {
		sendResponse(false, fmt.Sprintf("error creating directory: %v", err))
		return
	}

	sendResponse(true, "directory created successfully")
}
