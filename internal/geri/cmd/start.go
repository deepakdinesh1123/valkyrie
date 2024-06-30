package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/deepakdinesh1123/valkyrie/internal/geri/execute"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
	"github.com/deepakdinesh1123/valkyrie/internal/mq"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start geri",
	Long:  `Start geri`,
	RunE: func(cmd *cobra.Command, args []string) error {
		Start()
		return nil
	},
}

func Start() {
	fmt.Println("Starting geri...")
	ch, err := mq.GetChannel()
	if err != nil {
		logs.Logger.Err(err).Msg("Failed to get channel")
		panic(err)
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
		fmt.Printf("Failed to consume messages: %v\n", err)
		panic(err)
	}
	go func() {
		for execRequest := range execRequests {
			logs.Logger.Info().Msg(fmt.Sprintf("Received message: %s", string(execRequest.Body)))
			var executionRequest execution.ExecutionRequest
			err := json.Unmarshal(execRequest.Body, &executionRequest)
			if err != nil {
				logs.Logger.Err(err).Msg(fmt.Sprintf("Failed to unmarshal message: %s", string(execRequest.Body)))
				continue
			}
			res := execute.Execute(executionRequest)
			if res != "NotExecuted" {
				err := ch.Ack(execRequest.DeliveryTag, false)
				if err != nil {
					logs.Logger.Err(err).Msg(fmt.Sprintf("Failed to acknowledge message: %d", execRequest.DeliveryTag))
				}
			}
		}
	}()
	logs.Logger.Info().Msg("Worker started. Waiting for messages...")
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	<-sigterm
	logs.Logger.Info().Msg("Worker shutting down...")
}
