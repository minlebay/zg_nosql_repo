package shard_manager

import (
	cfg "go.uber.org/config"
)

type Config struct {
}

func NewManagerConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	if err := provider.Get("sharding").Populate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
