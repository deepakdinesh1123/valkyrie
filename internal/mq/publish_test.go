package mq_test

import (
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/mq"
)

func TestPublish(t *testing.T) {
	queue, err := mq.NewQueue("test", true, true, true, true, nil)
	if err != nil {
		t.Errorf("Failed to create queue: %v", err)
	}
	if queue == nil {
		t.Errorf("Queue is nil")
	}
	err = mq.Publish("test", []byte("test"))
	if err != nil {
		t.Errorf("Failed to publish message: %v", err)
	}
}
