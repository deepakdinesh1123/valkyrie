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
	var resp schemas.NewTerminalResponse

	if err := json.Unmarshal(data, &nt); err != nil {
		log.Printf("error parsing NewTerminal: %v", err)
		resp = schemas.NewTerminalResponse{Success: false, Msg: "Invalid request", TerminalID: ""}
		SendJSONMessage(ctx, c, resp)
		return
	}

	tty, tid, err := terminal.NewTTY(&nt)
	if err != nil {
		log.Printf("error creating terminal: %v", err)
		resp = schemas.NewTerminalResponse{Success: false, Msg: fmt.Sprintf("error creating terminal: %v", err)}
		SendJSONMessage(ctx, c, resp)
		return
	}

	s.AddTerminal(tid, *tty)
	resp = schemas.NewTerminalResponse{Success: true, Msg: "New Terminal created", TerminalID: tid}
	SendJSONMessage(ctx, c, resp)
}

func (s *Server) handleTerminalRead(ctx context.Context, c *websocket.Conn, data []byte) {
	var tr schemas.TerminalRead
	var resp schemas.TerminalReadResponse

	if err := json.Unmarshal(data, &tr); err != nil {
		log.Printf("error parsing TerminalRead: %v", err)
		resp = schemas.TerminalReadResponse{Success: false, Msg: "Invalid request", TerminalID: tr.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	tty, ok := s.GetTerminal(tr.TerminalID)
	if !ok {
		resp = schemas.TerminalReadResponse{Success: false, Msg: "Terminal not found", TerminalID: tr.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	output, err := tty.Read()
	if err != nil {
		resp = schemas.TerminalReadResponse{Success: false, Msg: fmt.Sprintf("Read error: %v", err), TerminalID: tr.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	resp = schemas.TerminalReadResponse{Success: true, Msg: "Read successful", Output: output, TerminalID: tr.TerminalID}
	SendJSONMessage(ctx, c, resp)
}

func (s *Server) handleTerminalWrite(ctx context.Context, c *websocket.Conn, data []byte) {
	var tw schemas.TerminalWrite
	var resp schemas.TerminalWriteResponse

	if err := json.Unmarshal(data, &tw); err != nil {
		log.Printf("error parsing TerminalWrite: %v", err)
		resp = schemas.TerminalWriteResponse{Success: false, Msg: "Invalid request", TerminalID: tw.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	tty, ok := s.GetTerminal(tw.TerminalID)
	if !ok {
		resp = schemas.TerminalWriteResponse{Success: false, Msg: "Terminal not found", TerminalID: tw.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	err := tty.Write([]byte(tw.Input))
	if err != nil {
		resp = schemas.TerminalWriteResponse{Success: false, Msg: fmt.Sprintf("Write error: %v", err), TerminalID: tw.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	resp = schemas.TerminalWriteResponse{Success: true, Msg: "Write successful", TerminalID: tw.TerminalID}
	SendJSONMessage(ctx, c, resp)
}

func (s *Server) handleTerminalClose(ctx context.Context, c *websocket.Conn, data []byte) {
	var tc schemas.TerminalClose
	var resp schemas.TerminalCloseResponse

	if err := json.Unmarshal(data, &tc); err != nil {
		log.Printf("error parsing TerminalClose: %v", err)
		resp = schemas.TerminalCloseResponse{Success: false, Msg: "Invalid request", TerminalID: tc.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	tty, ok := s.GetTerminal(tc.TerminalID)
	if !ok {
		resp = schemas.TerminalCloseResponse{Success: false, Msg: "Terminal not found", TerminalID: tc.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	err := tty.Close()
	if err != nil {
		resp = schemas.TerminalCloseResponse{Success: false, Msg: fmt.Sprintf("Close error: %v", err), TerminalID: tc.TerminalID}
		SendJSONMessage(ctx, c, resp)
		return
	}

	s.RemoveTerminal(tc.TerminalID)
	resp = schemas.TerminalCloseResponse{Success: true, Msg: "Terminal closed", TerminalID: tc.TerminalID}
	SendJSONMessage(ctx, c, resp)
}
