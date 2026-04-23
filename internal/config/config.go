package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectPath string `yaml:"project_path"`
	Port        int    `yaml:"port"`
	Watch       bool   `yaml:"watch"`
	Verbose     bool   `yaml:"verbose"`
	Workers     int    `yaml:"workers"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Default(), nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Workers == 0 {
		cfg.Workers = 4
	}
	return &cfg, nil
}

func Default() *Config {
	return &Config{
		ProjectPath: ".",
		Port:        8080,
		Watch:       false,
		Verbose:     false,
		Workers:     4,
	}
}

func FromEnv() *Config {
	return &Config{
		ProjectPath: getEnv("PROJECT_PATH", "."),
		Port:        getEnvInt("PORT", 8080),
		Watch:       getEnvBool("WATCH", false),
		Verbose:     getEnvBool("VERBOSE", false),
		Workers:     getEnvInt("WORKERS", 4),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err == nil {
			return i
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		return v == "true" || v == "1" || v == "yes"
	}
	return def
}