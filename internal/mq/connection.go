package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
)

var Connection *amqp.Connection

func GetConnection() (*amqp.Connection, error) {
	envConfig, err := config.GetEnvConfig()
	logger := logs.GetLogger()
	if err != nil {
		logger.Err(err).Msg("Failed to get env config")
		return nil, err
	}
	if Connection == nil {
		RABBITMQ_URL := fmt.Sprintf("amqp://guest:guest@%s:%s/", envConfig.RABBITMQ_HOST, envConfig.RABBITMQ_PORT)
		connection, err := amqp.Dial(
			RABBITMQ_URL,
		)
		if err != nil {
			logger.Err(err).Msg("Failed to connect to RabbitMQ")
			return nil, err
		}
		Connection = connection
	}
	return Connection, nil
}

func GetChannel() (*amqp.Channel, error) {
	logger := logs.GetLogger()
	connection, err := GetConnection()
	if err != nil {
		logger.Err(err).Msg("Failed to get connection")
		return nil, err
	}
	return connection.Channel()
}
