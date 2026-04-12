package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Device struct {
	MAC       string `yaml:"mac"`
	IP        string `yaml:"ip,omitempty"`
	Broadcast string `yaml:"broadcast,omitempty"`
	Password  string `yaml:"password,omitempty"`
}

type Config struct {
	Devices map[string]Device `yaml:"devices"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".wol.yml")
}

func Load() (*Config, error) {
	cfg := &Config{
		Devices: make(map[string]Device),
	}

	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.Devices == nil {
		cfg.Devices = make(map[string]Device)
	}
	return cfg, nil
}

func Save(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}
