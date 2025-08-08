package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port        string   `yaml:"port"`
	AllowedExts []string `yaml:"allowed_extensions"`
	MaxFiles    int      `yaml:"max_files_per_task"`
	MaxTasks    int      `yaml:"max_concurrent_tasks"`
}

func LoadConfig(configName string) (Config, error) {
	var newConfig Config

	data, err := os.ReadFile(configName)
	if err != nil {
		return newConfig, err
	}

	err = yaml.Unmarshal(data, &newConfig)
	return newConfig, err
}
