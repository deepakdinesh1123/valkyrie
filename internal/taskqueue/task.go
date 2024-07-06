package taskqueue

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Task interface {
	Execute(ctx context.Context) error
}

type TaskQueue struct {
	tasks   []Task
	Enqueue chan Task
	Dequeue chan Task
	Workers []*Worker
	logs    *zerolog.Logger
	mu      sync.Mutex
}

func NewTaskQueue(ctx context.Context, concurrency, bufferSize, taskTimeout int, logs *zerolog.Logger) *TaskQueue {

	tq := &TaskQueue{
		tasks:   make([]Task, 0),
		Enqueue: make(chan Task, bufferSize),
		Dequeue: make(chan Task),
		logs:    logs,
	}
	logs.Info().Msgf("Starting task queue with concurrency %d and buffer size %d", concurrency, bufferSize)
	for i := 0; i < concurrency; i++ {
		worker := NewWorker(i, tq, taskTimeout, logs)
		logs.Info().Msgf("Starting worker %d", i)
		worker.Start(ctx)
		tq.Workers = append(tq.Workers, worker)
	}
	return tq
}

func (tq *TaskQueue) Start() {
	go func() {
		for {
			select {
			case task := <-tq.Enqueue:
				tq.mu.Lock()
				tq.logs.Info().Msgf("Enqueued task: %v", task)
				tq.tasks = append(tq.tasks, task)
				tq.mu.Unlock()
			default:
				tq.mu.Lock()
				if len(tq.tasks) > 0 {
					tq.Dequeue <- tq.tasks[0]
					tq.logs.Info().Msgf("Dequeued task: %v", tq.tasks[0])
					tq.tasks = tq.tasks[1:]
				}
				tq.mu.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}
