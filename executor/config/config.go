package config

import (
	"github.com/spf13/viper"

	"github.com/deepakdinesh1123/valkyrie/executor/logger"
)

var log = logger.GetLogger()

type Environment struct {
	REDIS_HOST string `mapstructure:"REDIS_HOST"`
	REDIS_PORT string `mapstructure:"REDIS_PORT"`
}

var EnvConfig Environment

func init() {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env.dev")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Err(err)
	}
	err = viper.Unmarshal(&EnvConfig)
	if err != nil {
		log.Err(err)
	}
}
