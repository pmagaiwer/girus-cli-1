package common

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Language string `yaml:"language"`
}

var configPath string

// LoadConfig lê o arquivo de configuração padrão ~/.girus/config.yaml
func LoadConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		return &Config{Language: "pt"}
	}
	if configPath == "" {
		configPath = filepath.Join(home, ".girus", "config.yaml")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{Language: "pt"}
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &Config{Language: "pt"}
	}
	if cfg.Language == "" {
		cfg.Language = "pt"
	}
	return &cfg
}
