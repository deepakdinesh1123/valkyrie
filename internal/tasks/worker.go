package tasks

import (
	"github.com/google/uuid"
)

func GetNewWorker() error {
	worker := MachineryServer.NewWorker(uuid.NewString(), 10)

	return worker.Launch()
}
