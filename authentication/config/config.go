package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DatabaseUri  string
	DatabaseUser string
}

type ServerConfig struct {
	Port int
}

type Config struct {
	Db     *DBConfig
	Server *ServerConfig
}

func InitConfig() (*Config, error) {
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config *Config
	viper.Unmarshal(&config)
	if isValidConfig(config) {
		return nil, fmt.Errorf("Something wrong in Configuration. Please check configuration settings.")
	}
	return config, nil
}

func isValidConfig(config *Config) bool {
	return (config != nil &&
		config.Db != nil &&
		config.Db.DatabaseUri == "" &&
		config.Db.DatabaseUser == "")
}
