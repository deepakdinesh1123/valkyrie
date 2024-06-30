package mq

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue struct {
	Queue *amqp.Queue
}

var queues = map[string]*MessageQueue{}

func GetQueue(name string) (*MessageQueue, error) {
	if q, ok := queues[name]; ok {
		return q, nil
	}
	return nil, errors.New("queue not found")
}

func NewQueue(name string, Durable, AutoDelete, Exclusive, NoWait bool, Args map[string]interface{}) (*MessageQueue, error) {
	if _, ok := queues[name]; ok {
		return queues[name], nil
	}
	q := &MessageQueue{}

	amqpChannel, err := GetChannel()
	if err != nil {
		return nil, err
	}
	queue, err := amqpChannel.QueueDeclare(
		name,
		Durable,
		AutoDelete,
		Exclusive,
		NoWait,
		Args,
	)
	q.Queue = &queue
	if err != nil {
		return nil, err
	}
	queues[name] = q
	return q, nil
}
