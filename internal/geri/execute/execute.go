package execute

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

var (
	execDir = fmt.Sprintf("%s/%s", os.TempDir(), "valkyrie-execution/")
)

type ExecTask struct {
	executionRequest *execution.ExecutionRequest
	logs             *zerolog.Logger
	queries          *db.Queries
}

func NewExecTask(execRequest *execution.ExecutionRequest, queries *db.Queries, logs *zerolog.Logger) *ExecTask {
	return &ExecTask{
		executionRequest: execRequest,
		logs:             logs,
		queries:          queries,
	}
}

func (t *ExecTask) Execute(ctx context.Context) error {
	if _, err := os.Stat(execDir); os.IsNotExist(err) {
		err := os.MkdirAll(execDir, 0755)
		if err != nil {
			t.logs.Err(err).Msg("Failed to create execution directory")
			return err
		}
	}
	t.logs.Info().Msg("Execution directory already exists")
	executionPath := fmt.Sprintf("%s/%s/", execDir, t.executionRequest.ExecutionID)
	os.MkdirAll(executionPath, 0755)

	userCode := fmt.Sprintf("%s/%s", executionPath, t.executionRequest.File.Name)
	err := os.WriteFile(userCode, []byte(t.executionRequest.File.Content), 0744)
	if err != nil {
		t.logs.Err(err).Msg("Failed to write user code")
		return err
	}
	flakeNix := fmt.Sprintf("%s/%s", executionPath, "flake.nix")
	err = os.WriteFile(flakeNix, []byte(t.executionRequest.Environment), 0744)
	if err != nil {
		t.logs.Err(err).Msg("Failed to write flake.nix")
		return err
	}

	t.logs.Info().Msg("Running nix")
	cmd := exec.Command("nix", "develop")
	cmd.Dir = executionPath
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.logs.Err(err).Msgf("Failed to run nix: %s", string(err.Error()))
		t.queries.InsertExecutionResult(
			ctx,
			db.InsertExecutionResultParams{
				ExecutionID:     t.executionRequest.ExecutionID,
				Result:          pgtype.Text{String: string(out), Valid: true},
				ExecutionStatus: pgtype.Text{String: "Failed", Valid: true},
				ExecutedAt:      pgtype.Timestamp{Time: time.Now(), Valid: true},
			},
		)
		return err
	}
	fmt.Println(string(out))

	t.queries.InsertExecutionResult(
		ctx,
		db.InsertExecutionResultParams{
			ExecutionID:     t.executionRequest.ExecutionID,
			Result:          pgtype.Text{String: string(out), Valid: true},
			ExecutionStatus: pgtype.Text{String: "Executed", Valid: true},
			ExecutedAt:      pgtype.Timestamp{Time: time.Now(), Valid: true},
		},
	)
	return nil
}
