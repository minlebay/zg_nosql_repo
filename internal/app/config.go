package app

import (
	"go.uber.org/config"
	"go.uber.org/fx"
	"os"
)

type Config struct {
	Name string `yaml:"name"`
}

type ResultConfig struct {
	fx.Out
	Provider config.Provider
	Config   Config
}

func NewConfig() (ResultConfig, error) {
	yamlProvider, err := config.NewYAML(
		config.File("config.yaml"),
		config.Expand(os.LookupEnv),
	)
	if err != nil {
		return ResultConfig{}, err
	}

	cfg := Config{
		Name: "default",
	}

	if err = yamlProvider.Get("app").Populate(&cfg); err != nil {
		return ResultConfig{}, err
	}

	return ResultConfig{
		Provider: yamlProvider,
		Config:   cfg,
	}, nil
}
