package redis

import cfg "go.uber.org/config"

type Config struct {
	Address string `yaml:"address"`
	DB      string `yaml:"db"`
}

func NewRedisConfig(provider cfg.Provider) (*Config, error) {
	config := Config{}

	if err := provider.Get("redis").Populate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
