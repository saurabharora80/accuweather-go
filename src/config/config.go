package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"sync"
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

var (
	ConfigInstance      *Config
	ConfigInstanceError error
	onceForConfig       sync.Once
)

func GetConfig() (*Config, error) {
	onceForConfig.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("config")
		viper.AddConfigPath("../config")

		err := viper.ReadInConfig()

		if err != nil {
			fmt.Printf("Unable to read file %v\n", err)
			ConfigInstanceError = err
		}

		_ = viper.BindPFlags(pflag.CommandLine)
		viper.AutomaticEnv()

		config := Config{}

		err = viper.Unmarshal(&config)

		if err != nil {
			fmt.Printf("Unable to Unmarshall file %v\n", err)
			ConfigInstanceError = err
		} else {
			fmt.Println("==============Application Config================")
			fmt.Printf("%v\n", config)
			fmt.Println("===============================================")
			ConfigInstance = &config
		}
	})
	return ConfigInstance, ConfigInstanceError
}
