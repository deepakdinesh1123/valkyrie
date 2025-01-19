// server.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coder/websocket"

	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
)

type Message struct {
	MsgType string `json:"msgType"`
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Printf("Starting server on :1618")
	if err := http.ListenAndServe(":1618", nil); err != nil {
		log.Fatal(err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("websocket accept error: %v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	ctx := r.Context()

	for {
		_, data, err := c.Read(ctx)
		if err != nil {
			log.Printf("read error: %v", err)
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("json unmarshal error: %v", err)
			continue
		}

		switch msg.MsgType {
		case "Newterminal":
			{

			}
		case "Executecommand":
			var ec schemas.Executecommand
			if err := json.Unmarshal(data, &ec); err != nil {

			}
		case "Addfile":
			var af schemas.Addfile
			if err := json.Unmarshal(data, &af); err != nil {

			}
		}
	}
}
