package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/db"
	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/mq"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/database"
	"github.com/deepakdinesh1123/valkyrie/internal/taskqueue"
	"github.com/rs/zerolog"

	"github.com/deepakdinesh1123/valkyrie/internal/geri/execute"
	"github.com/jackc/pgx/v5"
)

type Consumer struct {
	DB           *pgx.Conn
	Queries      *db.Queries
	MessageQueue *mq.MessageQueue

	ValkyrieConfig *config.ValkyrieConfig
	EnvConfig      *config.Environment

	taskQueue *taskqueue.TaskQueue
}

func NewConsumer(ctx context.Context, logger *zerolog.Logger) (*Consumer, error) {
	DB, queries, err := database.GetDBConnection(ctx, logger)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to Postgres")
		return nil, err
	}
	valkyrieConfig, err := config.GetValkyrieConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to get valkyrie config")
	}
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		logger.Err(err).Msg("Failed to get env config")
		return nil, err
	}
	taskQueue := taskqueue.NewTaskQueue(ctx, valkyrieConfig.Geri.Concurrency, valkyrieConfig.Geri.BufferSize, valkyrieConfig.Geri.TaskTimeout, logger)
	taskQueue.Start()
	return &Consumer{
		DB:             DB,
		Queries:        queries,
		ValkyrieConfig: valkyrieConfig,
		taskQueue:      taskQueue,
		EnvConfig:      envConfig,
	}, nil
}

func (c *Consumer) Start(ctx context.Context, logger *zerolog.Logger) error {
	fmt.Println("Starting geri...")
	ch, err := mq.GetChannel()
	if err != nil {
		logger.Err(err).Msg("Failed to get channel")
		return err
	}
	execRequests, err := ch.Consume(
		"execute",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Err(err).Msg("Failed to consume messages")
		return err
	}
	go func() {
		for execRequest := range execRequests {
			logger.Info().Msg(fmt.Sprintf("Received message: %s", string(execRequest.Body)))
			var executionRequest execution.ExecutionRequest
			err := json.Unmarshal(execRequest.Body, &executionRequest)
			if err != nil {
				logger.Err(err).Msg(fmt.Sprintf("Failed to unmarshal message: %s", string(execRequest.Body)))
				continue
			}
			execTask := execute.NewExecTask(&executionRequest, c.Queries, logger)
			c.taskQueue.Enqueue <- execTask
		}
	}()
	logger.Info().Msg("Consumer started. Waiting for messages...")
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	<-sigterm
	logger.Info().Msg("Consumer shutting down...")

	return nil
}
