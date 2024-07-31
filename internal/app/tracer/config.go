package tracer

import (
	cfg "go.uber.org/config"
)

type Config struct {
	Enabled bool   `yaml:"enabled"`
	Url     string `yaml:"url"`
}

func NewTracerConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	if err := provider.Get("tracer").Populate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
