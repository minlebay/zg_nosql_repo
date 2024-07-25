package log

import "go.uber.org/config"

type Config struct {
	Url string `yaml:"url"`
}

func NewLogstashConfig(provider config.Provider) (*Config, error) {
	var config Config
	err := provider.Get("logstash").Populate(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
