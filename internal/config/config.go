package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
}

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type ServerConfig struct {
	Debug bool `yaml:"debug"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	path := fmt.Sprintf("configs/config.%s.yaml", env)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = "configs/config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Telegram.Token == "" {
		return nil, fmt.Errorf("telegram.token is not set in %s", path)
	}

	return &cfg, nil
}
