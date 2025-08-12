package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/agent/command"
	"github.com/deepakdinesh1123/valkyrie/agent/schemas"
	"github.com/deepakdinesh1123/valkyrie/agent/terminal"
	"github.com/rs/zerolog"
)

type Server struct {
	mu sync.RWMutex

	terminals map[string]terminal.TTY
	commands  map[string]command.Command
	logger    *zerolog.Logger
}

func NewServer(logger *zerolog.Logger) *Server {
	return &Server{
		terminals: make(map[string]terminal.TTY),
		commands:  make(map[string]command.Command),
		logger:    logger,
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
		s.logger.Err(err).Msg("websocket accept error")
		return
	}

	if c.Subprotocol() != "sandbox" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the sandbox subprotocol")
		return
	}

	defer c.Close(websocket.StatusInternalError, "")

	ctx := r.Context()

	for {
		_, data, err := c.Read(ctx)
		if err != nil {
			s.logger.Err(err).Msg("read error")
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			s.logger.Err(err).Msg("json unmarshal error occurred")
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
		case "ExecuteCommand":
			s.handleExecuteCommand(ctx, c, data)
		case "CommandReadOutput":
			s.handleCommandReadOutput(ctx, c, data)
		case "CommandWriteInput":
			s.handleCommandWriteInput(ctx, c, data)
		case "CommandTerminate":
			s.handleCommandTerminate(ctx, c, data)
		case "InstallNixPackage":
			s.handleInstallNixPackage(ctx, c, data)
		case "UninstallNixPackage":
			s.handleUninstallNixPackage(ctx, c, data)
		case "UpsertFile":
			s.handleUpsertFile(ctx, c, data)
		case "ReadFile":
			s.handleReadFile(ctx, c, data)
		case "DeleteFile":
			s.handleDeleteFile(ctx, c, data)
		case "UpsertDirectory":
			s.handleUpsertDirectory(ctx, c, data)
		case "ReadDirectory":
			s.handleReadDirectory(ctx, c, data)
		case "DeleteDirectory":
			s.handleDeleteDirectory(ctx, c, data)
		default:
			log.Printf("unknown message type: %v", msg.MsgType)

			response := schemas.Error{
				Message: fmt.Sprintf("Unrecognized message type: %s", msg.MsgType),
			}
			SendJSONMessage(ctx, c, response)
		}
	}
}

func (s *Server) Start(addr string) {
	http.HandleFunc("/sandbox", s.handleSandbox)

	s.logger.Info().Msgf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		s.logger.Fatal().Err(err)
	}
}
