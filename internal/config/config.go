package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const MicroserviceName = "Load GSC data"

type Config struct {
	Url     string `yaml:"url"`
	LogPath string `yaml:"log_path"`
	XmlPath string `yaml:"xml_path"`
	MaxProc uint64 `yaml:"max_proc"`
	Timeout uint64 `yaml:"timeout"`
}

func Setup(logger *zap.Logger) *Config {
	yamlCfg, err := loadYamlConfig()
	if err != nil {
		logger.Error("yaml config not found", zap.Error(err))
	}

	return yamlCfg
}

func loadYamlConfig() (*Config, error) {
	data, err := os.ReadFile("../../config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
