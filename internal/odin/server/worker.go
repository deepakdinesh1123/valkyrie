package server

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *OdinServer) DeleteExecutionWorker(ctx context.Context, params api.DeleteExecutionWorkerParams) (api.DeleteExecutionWorkerRes, error) {
	count, err := s.queries.WorkerTaskCount(ctx, pgtype.Int4{Int32: int32(params.WorkerId), Valid: true})
	if err != nil {
		return &api.DeleteExecutionWorkerBadRequest{
			Message: fmt.Sprintf("Failed to get worker: %v", err),
		}, nil
	}
	if !params.Force.Value && count > 0 {
		return &api.DeleteExecutionWorkerBadRequest{
			Message: fmt.Sprintf("Worker %d has %d tasks. Use force to delete", params.WorkerId, count),
		}, nil
	}
	err = s.queries.DeleteWorker(ctx, int32(params.WorkerId))
	if err != nil {
		return &api.DeleteExecutionWorkerInternalServerError{
			Message: fmt.Sprintf("Failed to delete worker: %v", err)}, nil
	}
	return &api.DeleteExecutionWorkerOK{}, nil
}
