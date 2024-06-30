package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(queue string, body []byte) error {
	amqpChannel, err := GetChannel()
	if err != nil {
		return err
	}
	err = amqpChannel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
