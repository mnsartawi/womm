package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Tool struct {
	Name             string            `yaml:"-"`
	Version          string            `yaml:"version"`
	InstallOverrides map[string]string `yaml:",inline"`
}

type Config struct {
	Dependencies map[string]Tool    `yaml:"dependencies"`
	Env          map[string]string  `yaml:"env"`
	Files        []string           `yaml:"files"`
	Commands     map[string][]string `yaml:"commands"`
	Platforms    map[string]bool    `yaml:"platforms"`

	Tools map[string]Tool `yaml:"-"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	for name, tool := range cfg.Dependencies {
		tool.Name = name
		cfg.Dependencies[name] = tool
	}

	if len(cfg.Dependencies) == 0 {
		rawMap := make(map[string]Tool)
		if err := yaml.Unmarshal(data, &rawMap); err == nil && len(rawMap) > 0 {
			cfg.Tools = rawMap
			for name, tool := range rawMap {
				tool.Name = name
				cfg.Tools[name] = tool
			}
		}
	}

	return cfg, nil
}

func (t Tool) GetCheckCmd() string {
	return t.Name + " --version"
}

func (t Tool) GetVersionRegex() string {
	return `\d+\.\d+\.\d+`
}

func (t Tool) GetPackageName() string {
	return t.Name
}
