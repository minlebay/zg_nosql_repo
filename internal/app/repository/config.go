package repository

import (
	cfg "go.uber.org/config"
)

type Config struct {
	Dbs []string `yaml:"mongodbs"`
}

func NewRepositoryConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	err := provider.Get("dbs").Populate(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
