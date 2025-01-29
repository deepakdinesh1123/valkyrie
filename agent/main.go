package main

import (
	"github.com/deepakdinesh1123/valkyrie/agent/server"
)

func main() {
	server := server.NewServer()
	server.Start(":1618")
}
