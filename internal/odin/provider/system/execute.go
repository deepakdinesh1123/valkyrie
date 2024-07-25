package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *SystemProvider) Execute(ctx context.Context, execReq db.Jobqueue) error {
	dir := fmt.Sprintf("%s/%s", s.envConfig.ODIN_SYSTEM_PROVIDER_BASE_DIR, execReq.CreatedAt.Time.Format("20060102150405"))

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err := os.WriteFile(fmt.Sprintf("%s/%s", dir, "flake.nix"), []byte(execReq.Flake.String), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s", dir, "main.py"), []byte(execReq.Script.String), os.ModePerm)
	if err != nil {
		return err
	}
	execCmd := exec.Command("nix", "run")

	execCmd.Dir = dir
	out, err := execCmd.CombinedOutput()
	s.queries.UpdateJob(ctx, db.UpdateJobParams{
		ID:   execReq.ID,
		Logs: pgtype.Text{String: string(out), Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}
