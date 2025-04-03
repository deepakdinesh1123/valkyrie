package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/deepakdinesh1123/valkyrie/pkg/api"
	"github.com/go-chi/chi/v5"
)

func (s *ValkyrieServer) CreateSandbox(ctx context.Context, req api.OptCreateSandbox, params api.CreateSandboxParams) (api.CreateSandboxRes, error) {
	if !s.envConfig.ENABLE_SANDBOX {
		return &api.CreateSandboxBadRequest{
			Message: "Sandbox functionality is not enabled",
		}, nil
	}
	sandboxId, err := s.sandboxService.AddSandbox(ctx, &req)
	if err != nil {
		s.logger.Err(err).Msg("could not create sandbox")
		return &api.CreateSandboxInternalServerError{
			Message: fmt.Sprintf("Failed to create sandbox: %v", err),
		}, nil
	}
	return &api.CreateSandboxOK{
		Result:           "Creating Sandbox....",
		SandboxId:        sandboxId,
		SandboxStatusSSE: api.NewOptString(fmt.Sprintf("/sandboxes/%d/status/sse", sandboxId)),
		SandboxStatusWS:  api.NewOptString(fmt.Sprintf("/sandboxes/%d/status/ws", sandboxId)),
	}, nil
}

func (s *ValkyrieServer) GetSandbox(ctx context.Context, params api.GetSandboxParams) (api.GetSandboxRes, error) {
	sandbox, err := s.queries.GetSandbox(ctx, params.SandboxId)
	if err != nil {
		return &api.Error{
			Message: fmt.Sprintf("error getting sandbox %s", err),
		}, nil
	}
	if sandbox.CurrentState == "pending" || sandbox.CurrentState == "creating" {
		return &api.GetSandboxOK{
			Type: api.SandboxStateGetSandboxOK,
			SandboxState: api.SandboxState{
				State:     sandbox.CurrentState,
				SandboxId: params.SandboxId,
			},
		}, nil
	} else if sandbox.CurrentState == "failed" {
		return &api.Error{
			Message: sandbox.Details.Error,
		}, nil
	} else {
		return &api.GetSandboxOK{
			Type: api.SandboxGetSandboxOK,
			Sandbox: api.Sandbox{
				SandboxId: sandbox.SandboxID,
				State:     sandbox.CurrentState,
				URL:       sandbox.SandboxUrl.String,
				CreatedAt: sandbox.CreatedAt.Time,
				AgentURL:  sandbox.SandboxAgentUrl.String,
			},
		}, nil
	}
}

func (s *ValkyrieServer) GetSandboxSSE(w http.ResponseWriter, r *http.Request) {
	// Parse and validate sandbox ID
	sandboxIDStr := chi.URLParam(r, "sandboxId")
	if sandboxIDStr == "" {
		s.logger.Err(nil).Msg("missing sandboxId")
		http.Error(w, "missing sandboxId", http.StatusBadRequest)
		return
	}

	sandboxID, err := strconv.ParseInt(sandboxIDStr, 10, 64)
	if err != nil {
		s.logger.Err(err).Msg("error parsing sandbox_id")
		http.Error(w, fmt.Sprintf("error parsing sandbox_id: %v", err), http.StatusBadRequest)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := r.Context()
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Use a more reasonable ticker interval
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Helper function to send SSE messages
	sendMessage := func(eventType string, data interface{}) error {
		msg, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshaling message: %v", err)
		}

		if eventType != "" {
			fmt.Fprintf(w, "event: %s\n", eventType)
		}
		fmt.Fprintf(w, "data: %s\n\n", string(msg))
		flusher.Flush()
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			// Send a final message before closing
			sendMessage("close", api.Error{Message: "stream closed by client"})
			return
		case <-ticker.C:
			sandbox, err := s.queries.GetSandbox(ctx, sandboxID)
			if err != nil {
				sendMessage("error", api.Error{
					Message: fmt.Sprintf("error fetching sandbox: %v", err),
				})
				return
			}

			// Handle in-progress states
			if sandbox.CurrentState == "pending" || sandbox.CurrentState == "creating" {
				err = sendMessage("status", api.Sandbox{
					State:     sandbox.CurrentState,
					SandboxId: sandboxID,
				})
				if err != nil {
					sendMessage("error", api.Error{Message: err.Error()})
					return
				}
				continue
			}

			// Handle failed state
			if sandbox.CurrentState == "failed" {
				err = sendMessage("error", api.Error{
					Message: sandbox.Details.Message,
				})
				if err != nil {
					sendMessage("error", api.Error{Message: err.Error()})
				}
				return // Exit stream on failure
			}

			err = sendMessage("status", api.Sandbox{
				SandboxId: sandbox.SandboxID,
				State:     sandbox.CurrentState,
				URL:       sandbox.SandboxUrl.String,
				CreatedAt: sandbox.CreatedAt.Time,
				AgentURL:  sandbox.SandboxAgentUrl.String,
			})
			if err != nil {
				sendMessage("error", api.Error{Message: err.Error()})
				return
			}
		}
	}
}

func (s *ValkyrieServer) GetSandboxWS(w http.ResponseWriter, r *http.Request) {
	// Parse and validate sandbox ID
	sandboxIDStr := chi.URLParam(r, "sandboxId")
	if sandboxIDStr == "" {
		s.logger.Err(nil).Msg("missing sandboxId")
		http.Error(w, "missing sandboxId", http.StatusBadRequest)
		return
	}

	sandboxID, err := strconv.ParseInt(sandboxIDStr, 10, 64)
	if err != nil {
		s.logger.Err(err).Msg("error parsing sandbox_id")
		http.Error(w, fmt.Sprintf("error parsing sandbox_id: %v", err), http.StatusBadRequest)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Adjust based on security needs
	})
	if err != nil {
		s.logger.Err(err).Msg("error establishing websocket connection")
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing connection")

	ctx := r.Context()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	sendMessage := func(eventType string, data interface{}) error {
		msg, err := json.Marshal(map[string]interface{}{
			"event": eventType,
			"data":  data,
		})
		if err != nil {
			return fmt.Errorf("error marshaling message: %v", err)
		}
		return conn.Write(ctx, websocket.MessageText, msg)
	}

	for {
		select {
		case <-ctx.Done():
			sendMessage("close", api.Error{Message: "stream closed by client"})
			return
		case <-ticker.C:
			sandbox, err := s.queries.GetSandbox(ctx, sandboxID)
			if err != nil {
				sendMessage("error", api.Error{Message: fmt.Sprintf("error fetching sandbox: %v", err)})
				return
			}

			if sandbox.CurrentState == "pending" || sandbox.CurrentState == "creating" {
				err = sendMessage("status", api.Sandbox{
					State:     sandbox.CurrentState,
					SandboxId: sandboxID,
				})
				if err != nil {
					sendMessage("error", api.Error{Message: err.Error()})
					return
				}
				continue
			}

			if sandbox.CurrentState == "failed" {
				sendMessage("error", api.Error{Message: sandbox.Details.Message})
				return
			}

			if sandbox.CurrentState == "running" {
				err = sendMessage("status", api.Sandbox{
					SandboxId: sandbox.SandboxID,
					State:     sandbox.CurrentState,
					URL:       sandbox.SandboxUrl.String,
					CreatedAt: sandbox.CreatedAt.Time,
					AgentURL:  sandbox.SandboxAgentUrl.String,
				})
				if err != nil {
					sendMessage("error", api.Error{Message: err.Error()})
					return
				}
				return
			}
		}
	}
}
