package tasks

import (
	"fmt"

	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"

	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"

	ValkyrieConfig "github.com/deepakdinesh1123/valkyrie/internal/config"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
)

var MachineryServer *machinery.Server

func init() {
	REDIS_URL := fmt.Sprintf("%s:%s", ValkyrieConfig.EnvConfig.REDIS_HOST, ValkyrieConfig.EnvConfig.REDIS_PORT)
	if MachineryServer != nil {
		return
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
	broker := redisbroker.NewGR(cnf, []string{REDIS_URL}, 0)
	backend := redisbackend.NewGR(cnf, []string{REDIS_URL}, 1)
	lock := eagerlock.New()
	MachineryServer = machinery.NewServer(cnf, broker, backend, lock)
	err := MachineryServer.RegisterTasks(TasksMap)
	logs.Logger.Info().Msg("Machinery Server started and tasks registered")
	if err != nil {
		logs.Logger.Fatal().Msg("Could not register tasks")
	}
}
