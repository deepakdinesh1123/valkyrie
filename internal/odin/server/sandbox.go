package server

import (
	"context"
	"fmt"
	"net/url"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/db/jsonschema"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

func (s *OdinServer) CreateSandbox(ctx context.Context, params api.CreateSandboxParams) (api.CreateSandboxRes, error) {
	res, err := s.queries.AddSandboxJobTx(ctx, db.AddSandboxTxParams{
		SandboxConfig: jsonschema.SandboxConfig{},
		MaxRetries:    s.envConfig.ODIN_MAX_RETRIES,
	})
	if err != nil {
		s.logger.Err(err).Msg("could not create sandbox")
		return &api.CreateSandboxOK{
			Result: "Failure:",
		}, nil
	}
	return &api.CreateSandboxOK{
		Result:    "Creating Sandbox",
		SandboxId: res.SandboxId,
	}, nil
}

func (s *OdinServer) GetSandbox(ctx context.Context, params api.GetSandboxParams) (api.GetSandboxRes, error) {
	sandbox, err := s.queries.GetSandbox(ctx, params.SandboxId)
	if err != nil {
		return &api.Error{
			Message: fmt.Sprintf("error getting sandbox %s", err),
		}, nil
	}
	if sandbox.CurrentState == "pending" || sandbox.CurrentState == "creating" {
		return &api.Sandbox{
			State:     sandbox.CurrentState,
			SandboxId: params.SandboxId,
		}, nil
	} else {
		sandboxURL, err := url.Parse(sandbox.SandboxUrl.String)
		if err != nil {
			return &api.Error{
				Message: fmt.Sprintf("error parsing url %s", err),
			}, nil
		}
		return &api.Sandbox{
			SandboxId: sandbox.SandboxID,
			State:     sandbox.CurrentState,
			URL:       api.NewOptURI(*sandboxURL),
			CreatedAt: api.NewOptDateTime(sandbox.CreatedAt.Time),
		}, nil
	}
}
