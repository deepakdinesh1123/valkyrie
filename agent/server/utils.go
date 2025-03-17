package server

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
)

// SendJSONMessage marshals a struct and sends it over a WebSocket connection.
func SendJSONMessage(ctx context.Context, conn *websocket.Conn, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		return err
	}

	return nil
}
