package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Bot struct {
		TelegramToken string `yaml:"telegram_token"`
	}

	DB struct {
		DBHost     string `yaml:"host"`
		DBPort     string `yaml:"port"`
		DBUser     string `yaml:"user"`
		DBPassword string `yaml:"password"`
		DBName     string `yaml:"name"`
	}
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil

}
