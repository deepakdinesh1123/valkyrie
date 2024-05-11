package rust

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
)

const (
	TempDir    = "/skald/rust/"
	MainRsPath = "src/main.rs"
)

func Execute(ctx context.Context, executionRequest string) error {
	var execRequest execution.ExecutionRequest
	_ = json.Unmarshal([]byte(executionRequest), &execRequest)

	fmt.Println(execRequest, executionRequest, execRequest.SystemDependencies)
	return nil
}
