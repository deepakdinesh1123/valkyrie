package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/agent/command"
	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
)

func (s *Server) AddCommand(id string, cmd command.Command) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commands[id] = cmd
}

func (s *Server) RemoveCommand(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.commands, id)
}

func (s *Server) GetCommand(id string) (command.Command, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.commands[id]
	return c, ok
}

func (s *Server) handleExecuteCommand(ctx context.Context, c *websocket.Conn, data []byte) {
	var ec schemas.ExecuteCommand
	msgType := "ExecuteCommandResponse"

	// Helper function to send a response
	sendResponse := func(state schemas.State, cid string, success bool, msg string) {
		resp := schemas.ExecuteCommandResponse{
			State:     &state,
			MsgType:   &msgType,
			CommandID: cid,
			Success:   success,
			Msg:       msg,
		}
		SendJSONMessage(ctx, c, resp)
	}

	// Parse request
	if err := json.Unmarshal(data, &ec); err != nil {
		log.Printf("error parsing ExecuteCommand: %v", err)
		sendResponse(schemas.Exited, "", false, "Failed to parse ExecuteCommand")
		return
	}

	// Create and add command
	cmd, cid, err := command.NewCommand(&ec)
	if err != nil {
		sendResponse(schemas.Exited, "", false, "Failed to create command")
		return
	}

	s.AddCommand(cid, *cmd)
	sendResponse(schemas.Running, cid, true, "Command is running")

	go func(cmd *command.Command) {
		err = cmd.Wait()
		if err != nil {
			log.Printf("error when waiting for command: %v\n", err)
			sendResponse(schemas.Exited, cid, false, "Command execution failed")
		} else {
			sendResponse(schemas.Exited, cid, true, "Command execution completed")
		}
		cmd.Completed = true
	}(cmd)
}

func (s *Server) handleCommandWriteInput(ctx context.Context, c *websocket.Conn, data []byte) {
	var cwi schemas.CommandWriteInput
	msgType := "CommandWriteInputResponse"

	sendResponse := func(success bool, commandID string) {
		SendJSONMessage(ctx, c, schemas.CommandWriteInputResponse{
			MsgType:   &msgType,
			Success:   success,
			CommandID: commandID,
		})
	}

	if err := json.Unmarshal(data, &cwi); err != nil {
		log.Printf("error parsing CommandWriteInput: %v", err)
		sendResponse(false, cwi.CommandID)
		return
	}

	cmd, ok := s.GetCommand(cwi.CommandID)
	if !ok {
		sendResponse(false, cwi.CommandID)
		return
	}

	if err := cmd.Write([]byte(*cwi.Input)); err != nil {
		sendResponse(false, cwi.CommandID)
		return
	}

	sendResponse(true, cwi.CommandID)
}

func (s *Server) handleCommandReadOutput(ctx context.Context, c *websocket.Conn, data []byte) {
	var cro schemas.CommandReadOutput
	msgType := "CommandReadOutputResponse"

	sendResponse := func(output string, commandID string) {
		SendJSONMessage(ctx, c, schemas.CommandReadOutputResponse{
			MsgType:   &msgType,
			Stdout:    output,
			CommandID: commandID,
		})
	}

	if err := json.Unmarshal(data, &cro); err != nil {
		log.Printf("error parsing CommandReadOutput: %v", err)
		sendResponse("", cro.CommandID)
		return
	}

	cmd, ok := s.GetCommand(cro.CommandID)
	if !ok {
		sendResponse("", cro.CommandID)
		return
	}

	output, err := cmd.Read()
	if err != nil {
		log.Printf("error reading from output: %v", err)
		sendResponse("", cro.CommandID)
		return
	}

	sendResponse(output, cro.CommandID)

	if cmd.Completed {
		s.RemoveCommand(cro.CommandID)
	}
}

func (s *Server) handleCommandTerminate(ctx context.Context, c *websocket.Conn, data []byte) {
	var ct schemas.CommandTerminate
	msgType := "CommandTerminateResponse"

	sendResponse := func(commandID string) {
		SendJSONMessage(ctx, c, schemas.CommandTerminateResponse{
			MsgType:   &msgType,
			CommandID: commandID,
		})
	}

	if err := json.Unmarshal(data, &ct); err != nil {
		log.Printf("error parsing CommandTerminate: %v", err)
		sendResponse(ct.CommandID)
		return
	}

	cmd, ok := s.GetCommand(ct.CommandID)
	if !ok {
		sendResponse(ct.CommandID)
		return
	}

	if err := cmd.Terminate(); err != nil {
		sendResponse(ct.CommandID)
		return
	}

	s.RemoveCommand(ct.CommandID)
	sendResponse(ct.CommandID)
}
