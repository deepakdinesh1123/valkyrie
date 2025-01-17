//go:build system || all || darwin

package system

import (
	"context"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor/system/native"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/executor/system/nsjail"
)

func getSystemExecutorClient(ctx context.Context, se *SystemExecutor) (SystemExecutionClient, error) {
	fmt.Println("System executor client is ........................", se.envConfig.ODIN_WORKER_SYSTEM_EXECUTOR)
	switch se.envConfig.ODIN_WORKER_SYSTEM_EXECUTOR {
	case "nsjail":
		return nsjail.NewNSJailExecutor(
			ctx,
			se.envConfig,
			se.queries,
			se.workerId,
			se.tp,
			se.mp,
			se.logger,
		)
	case "native":
		return native.NewSystemExecutor(
			ctx,
			se.envConfig,
			se.queries,
			se.workerId, se.tp,
			se.mp,
			se.logger,
		)
	default:
		return nil, fmt.Errorf("No client with given name")
	}
}
