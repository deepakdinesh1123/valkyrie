package server

import (
	"github.com/deepakdinesh1123/valkyrie/internal/mq"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

type Server struct{}

func NewServer() *api.Server {
	createQueues()
	server := &Server{}
	srv, err := api.NewServer(server)
	if err != nil {
		panic(err)
	}
	return srv
}

func createQueues() error {
	_, err := mq.NewQueue("execute", true, true, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}
