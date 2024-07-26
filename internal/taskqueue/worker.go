package taskqueue

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Worker struct {
	id          int
	taskQueue   *TaskQueue
	taskTimeout int

	logs *zerolog.Logger
}

func NewWorker(id int, taskQueue *TaskQueue, taskTimeout int, logs *zerolog.Logger) *Worker {
	return &Worker{id: id, taskQueue: taskQueue, logs: logs, taskTimeout: taskTimeout}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			task := <-w.taskQueue.Dequeue
			err := ExecuteTask(ctx, task, w.taskTimeout)
			if err != nil {
				w.logs.Err(err).Msgf("Worker %d: failed to execute task", w.id)
			} else {
				w.logs.Info().Msgf("Worker %d: executed task", w.id)
			}
		}
	}()
}

func ExecuteTask(ctx context.Context, task Task, taskTimeout int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(taskTimeout)*time.Second)
	defer cancel()
	return task.Execute(ctx)
}
