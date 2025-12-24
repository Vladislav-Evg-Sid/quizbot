package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Service  ServiceConfig  `yaml:"service"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type ServiceConfig struct {
	GRPCPort    int    `yaml:"grpc_port"`
	RESTPort    int    `yaml:"rest_port"`
	SwaggerPath string `yaml:"swagger_path"`
}
type KafkaConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	TopicNameAddQuizResult string `yaml:"topic_name_add_quiz_results"`
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
