package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver    string `mapstructure:"DB_DRIVER"`
	DBSource    string `mapstructure:"DB_SOURCE"`
	HTTPAddress string `mapstructure:"HTTP_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
