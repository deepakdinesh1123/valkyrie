package main

import (
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/tasks"
)

func main() {
	err := tasks.GetNewWorker()
	if err != nil {
		logs.Logger.Fatal().Msg("Could not create new worker")
	}
	logs.Logger.Info().Msg("Worker started")
}
