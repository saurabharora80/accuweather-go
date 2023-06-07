package main

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Upstream struct {
		Host                         string        `mapstructure:"host"`
		MaxIdleConnections           int           `mapstructure:"max-idle-connections"`
		IdleConnectionTimeoutSeconds time.Duration `mapstructure:"idle-connection-timeout-seconds"`
		Key                          string        `mapstructure:"key"`
	} `mapstructure:"upstream"`
}

func config() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resources")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("Unable to read file %v\n", err)
		return Config{}, err
	}

	config := Config{}

	err = viper.Unmarshal(&config)

	if err != nil {
		fmt.Printf("Unable to Unmarshall file %v\n", err)
		return Config{}, err
	} else {
		return config, err
	}

}
