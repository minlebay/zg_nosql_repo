package repository

import (
	cfg "go.uber.org/config"
)

type Config struct {
	Url          string `yaml:"url"`
	DatabaseName string `yaml:"database_name"`
}

func NewRepositoryConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	if err := provider.Get("mongodb").Populate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
