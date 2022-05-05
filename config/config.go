package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BaseUrl string `mapstructure:"base_url"`
}

func (c Config) Load(env string) (Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType(env)
	viper.AddConfigPath("/app/config")

	viper.AutomaticEnv()

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
