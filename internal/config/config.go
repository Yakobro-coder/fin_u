package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DB     DB         `mapstructure:"db"`
	HTTP   HTTPConfig `mapstructure:"app"`
	Logger Logger     `mapstructure:"app"`
}

type DB struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	DBName   string `mapstructure:"POSTGRES_DB"`
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
}

type HTTPConfig struct {
	Port int `mapstructure:"HTTP_PORT"`
}

type Logger struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

func Parse() (*Config, error) {
	var cfg Config

	//serverEnvironment := os.Getenv("SERVER_ENV")
	serverEnvironment := "dev"

	viper.SetConfigName(serverEnvironment)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./env/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to parse config from %s.yaml: %w", serverEnvironment, err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from %s.yaml: %w", serverEnvironment, err)
	}

	return &cfg, nil
}
