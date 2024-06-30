package mq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
)

var Connection *amqp.Connection

func GetConnection() (*amqp.Connection, error) {
	if Connection == nil {
		RABBITMQ_URL := fmt.Sprintf("amqp://guest:guest@%s:%s/", config.EnvConfig.RABBITMQ_HOST, config.EnvConfig.RABBITMQ_PORT)
		connection, err := amqp.Dial(
			RABBITMQ_URL,
		)
		if err != nil {
			logs.Logger.Err(err).Msg("Failed to connect to RabbitMQ")
			return nil, err
		}
		Connection = connection
	}
	return Connection, nil
}

func GetChannel() (*amqp.Channel, error) {
	if Connection == nil {
		connection, err := GetConnection()
		if err != nil {
			logs.Logger.Err(err).Msg("Failed to get connection")
			return nil, err
		}
		return connection.Channel()
	}
	return Connection.Channel()
}
