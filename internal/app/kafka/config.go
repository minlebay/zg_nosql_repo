package kafka

import (
	cfg "go.uber.org/config"
)

type Config struct {
	Address  string `yaml:"address"`
	GroupID  string `yaml:"group_id"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Topics   string `yaml:"topic"`
}

func NewKafkaConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	if err := provider.Get("kafka").Populate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
