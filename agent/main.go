package main

import (
	"github.com/deepakdinesh1123/valkyrie/agent/logs"
	"github.com/deepakdinesh1123/valkyrie/agent/server"
)

func main() {
	logger := logs.GetLogger()

	server := server.NewServer(logger)
	server.Start(":1618")
}
