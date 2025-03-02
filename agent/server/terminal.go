package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
	"github.com/deepakdinesh1123/valkyrie/agent/terminal"
)

func (s *Server) AddTerminal(id string, t terminal.TTY) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.terminals[id] = t
}

func (s *Server) RemoveTerminal(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.terminals, id)
}

func (s *Server) GetTerminal(id string) (terminal.TTY, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.terminals[id]
	return t, ok
}

func (s *Server) handleNewTerminal(ctx context.Context, c *websocket.Conn, data []byte) {
	var nt schemas.NewTerminal
	msgType := "NewTerminalResponse"

	sendResponse := func(id string, success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.NewTerminalResponse{
			MsgType:    &msgType,
			Success:    success,
			Msg:        msg,
			TerminalID: id,
		})
	}

	if err := json.Unmarshal(data, &nt); err != nil {
		log.Printf("error parsing NewTerminal: %v", err)
		sendResponse("", false, "Invalid request")
		return
	}

	tty, tid, err := terminal.NewTTY(&nt)
	if err != nil {
		log.Printf("error creating terminal: %v", err)
		sendResponse("", false, fmt.Sprintf("error creating terminal: %v", err))
		return
	}

	s.AddTerminal(tid, *tty)
	sendResponse(tid, true, "New Terminal created")
}

func (s *Server) handleTerminalRead(ctx context.Context, c *websocket.Conn, data []byte) {
	var tr schemas.TerminalRead

	msgType := "TerminalReadResponse"

	sendResponse := func(id string, success bool, msg, output string) {
		SendJSONMessage(ctx, c, schemas.TerminalReadResponse{
			MsgType:    &msgType,
			Success:    success,
			Msg:        msg,
			Output:     output,
			TerminalID: id,
		})
	}

	if err := json.Unmarshal(data, &tr); err != nil {
		log.Printf("error parsing TerminalRead: %v", err)
		sendResponse(tr.TerminalID, false, "Invalid request", "")
		return
	}

	tty, ok := s.GetTerminal(tr.TerminalID)
	if !ok {
		sendResponse(tr.TerminalID, false, "Terminal not found", "")
		return
	}

	output, err := tty.Read()
	if err != nil {
		sendResponse(tr.TerminalID, false, fmt.Sprintf("Read error: %v", err), "")
		return
	}

	sendResponse(tr.TerminalID, true, "Read successful", output)
}

func (s *Server) handleTerminalWrite(ctx context.Context, c *websocket.Conn, data []byte) {
	var tw schemas.TerminalWrite

	msgType := "TerminalWriteResponse"

	sendResponse := func(id string, success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.TerminalWriteResponse{
			MsgType:    &msgType,
			Success:    success,
			Msg:        msg,
			TerminalID: id,
		})
	}

	if err := json.Unmarshal(data, &tw); err != nil {
		log.Printf("error parsing TerminalWrite: %v", err)
		sendResponse(tw.TerminalID, false, "Invalid request")
		return
	}

	tty, ok := s.GetTerminal(tw.TerminalID)
	if !ok {
		sendResponse(tw.TerminalID, false, "Terminal not found")
		return
	}

	if err := tty.Write([]byte(tw.Input)); err != nil {
		sendResponse(tw.TerminalID, false, fmt.Sprintf("Write error: %v", err))
		return
	}

	sendResponse(tw.TerminalID, true, "Write successful")
}

func (s *Server) handleTerminalClose(ctx context.Context, c *websocket.Conn, data []byte) {
	var tc schemas.TerminalClose

	msgType := "TerminalReadResponse"

	sendResponse := func(id string, success bool, msg string) {
		SendJSONMessage(ctx, c, schemas.TerminalCloseResponse{
			MsgType:    &msgType,
			Success:    success,
			Msg:        msg,
			TerminalID: id,
		})
	}

	if err := json.Unmarshal(data, &tc); err != nil {
		log.Printf("error parsing TerminalClose: %v", err)
		sendResponse(tc.TerminalID, false, "Invalid request")
		return
	}

	tty, ok := s.GetTerminal(tc.TerminalID)
	if !ok {
		sendResponse(tc.TerminalID, false, "Terminal not found")
		return
	}

	if err := tty.Close(); err != nil {
		sendResponse(tc.TerminalID, false, fmt.Sprintf("Close error: %v", err))
		return
	}

	s.RemoveTerminal(tc.TerminalID)
	sendResponse(tc.TerminalID, true, "Terminal closed")
}
