package database_test

import (
	"context"
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/database"
)

func TestGetDBConnection(t *testing.T) {
	_, _, err := database.GetDBConnection(context.Background(), nil)
	if err != nil {
		t.Errorf("Failed to connect to Postgres: %v", err)
	}
}
