package tasks

import (
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"

	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
)

var MachineryServer *machinery.Server

func GetMachineryServer() (*machinery.Server, error) {
	if MachineryServer != nil {
		return MachineryServer, nil
	}
	cnf := &config.Config{
		DefaultQueue: "machinery_tasks",
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}
	broker := redisbroker.NewGR(cnf, []string{"localhost:6379"}, 0)
	backend := redisbackend.NewGR(cnf, []string{"localhost:6379"}, 1)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)
	return server, server.RegisterTasks(
		TasksMap,
	)
}
