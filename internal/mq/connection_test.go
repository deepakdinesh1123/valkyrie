// go:build integration

package mq_test

import (
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/mq"
)

func TestGetConnection(t *testing.T) {
	conn, err := mq.GetConnection()
	if err != nil {
		t.Errorf("Failed to get connection: %v", err)
	}
	if conn == nil {
		t.Errorf("Connection is nil")
	}
}

func TestGetChannel(t *testing.T) {
	ch, err := mq.GetChannel()
	if err != nil {
		t.Errorf("Failed to get channel: %v", err)
	}
	if ch == nil {
		t.Errorf("Channel is nil")
	}
}
