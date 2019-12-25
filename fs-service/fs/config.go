package fs

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	User     string
	Password string
	DBName   string
}

type ClientEndpoints struct {
	Endpoints []string
}

type Config struct {
	Db     DBConfig
	Client ClientEndpoints
}

func InitConfig() (interface{}, error) {
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	viper.Unmarshal(&config)
	return config, nil
}
