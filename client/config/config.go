package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Redis   RedisConfig   `yaml:"redis"`
	TgBot   TgBotConfig   `yaml:"tgbot"`
	Network NetworkConfig `yaml:"network"`
}

type RedisConfig struct {
	Url      string `yaml:"url"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type TgBotConfig struct {
	Token string `yaml:"token"`
}

type NetworkConfig struct {
	AdminREST  string `yaml:"admin_rest"`
	AdminGRPC  string `yaml:"admin_grpc"`
	PlayerREST string `yaml:"player_rest"`
	PlayerGRPC string `yaml:"player_grpc"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	expandedData := os.ExpandEnv(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
