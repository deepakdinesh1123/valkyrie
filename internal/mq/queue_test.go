package mq_test

import (
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/mq"
)

func TestNewQueue(t *testing.T) {
	queue, err := mq.NewQueue("test", true, true, true, true, nil)
	if err != nil {
		t.Errorf("Failed to create queue: %v", err)
	}
	if queue == nil {
		t.Errorf("Queue is nil")
	}
}
