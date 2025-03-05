package server

import (
	"context"

	"github.com/coder/websocket"
)

func (s *Server) handleDeleteDirectory(ctx context.Context, c *websocket.Conn, data []byte) {}
func (s *Server) handleReadDirectory(ctx context.Context, c *websocket.Conn, data []byte)   {}
func (s *Server) handleUpsertDirectory(ctx context.Context, c *websocket.Conn, data []byte) {}
