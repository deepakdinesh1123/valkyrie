package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/agent/command"
	"github.com/deepakdinesh1123/valkyrie/agent/terminal"
)

type Server struct {
	mu sync.RWMutex

	terminals map[string]terminal.TTY
	commands  map[string]command.Command
}

func NewServer() *Server {
	return &Server{
		terminals: make(map[string]terminal.TTY),
		commands:  make(map[string]command.Command),
	}
}

type Message struct {
	MsgType string `json:"msgType"`
}

func (s *Server) handleSandbox(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		Subprotocols:       []string{"sandbox"},
	})
	if err != nil {
		log.Printf("websocket accept error: %v", err)
		return
	}

	if c.Subprotocol() != "sandbox" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the terminal subprotocol")
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
		case "NewTerminal":
			s.handleNewTerminal(ctx, c, data)
		case "TerminalRead":
			s.handleTerminalRead(ctx, c, data)
		case "TerminalWrite":
			s.handleTerminalWrite(ctx, c, data)
		case "TerminalClose":
			s.handleTerminalClose(ctx, c, data)
		default:
			log.Printf("unknown message type: %v", msg.MsgType)
		}
	}
}

func (s *Server) Start(addr string) {
	http.HandleFunc("/sandbox", s.handleSandbox)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
