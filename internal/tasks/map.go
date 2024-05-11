package tasks

import (
	"github.com/deepakdinesh1123/valkyrie/internal/skald/environments/rust"
)

var TasksMap = map[string]interface{}{
	"execute_rust": rust.Execute,
}
