package nosql

import (
	cfg "go.uber.org/config"
)

type Config struct {
}

func NewRepoConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	return &config, nil
}
