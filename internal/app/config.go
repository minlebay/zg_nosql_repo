package app

import (
	"fmt"
	"go.uber.org/config"
	"go.uber.org/fx"
	"os"
	"strings"
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
	yamlProvider, err := config.NewYAML(config.File("config.yaml"))
	if err != nil {
		return ResultConfig{}, err
	}

	var yamlConfig map[interface{}]interface{}
	if err = yamlProvider.Get(config.Root).Populate(&yamlConfig); err != nil {
		return ResultConfig{}, err
	}

	stringMap := convertMapKeysToStrings(yamlConfig)
	replaceEnvVariables(stringMap)

	provider, err := config.NewYAML(config.Static(stringMap))
	if err != nil {
		return ResultConfig{}, err
	}

	config := Config{
		Name: "default",
	}

	if err = provider.Get("app").Populate(&config); err != nil {
		return ResultConfig{}, err
	}

	return ResultConfig{
		Provider: provider,
		Config:   config,
	}, nil
}

func convertMapKeysToStrings(m map[interface{}]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		strKey := fmt.Sprintf("%v", k)
		switch val := v.(type) {
		case map[interface{}]interface{}:
			newMap[strKey] = convertMapKeysToStrings(val)
		default:
			newMap[strKey] = val
		}
	}
	return newMap
}

func replaceEnvVariables(config map[string]interface{}) {
	for key, value := range config {
		switch v := value.(type) {
		case map[string]interface{}:
			replaceEnvVariables(v)
		case string:
			if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
				envVar := v[2 : len(v)-1]
				config[key] = os.Getenv(envVar)
			}
		}
	}
}
