package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Service struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"service"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"postgres"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	Auth struct {
		RefreshLifeTime int64  `yaml:"refresh_life_time"`
		AccessLifeTime  int64  `yaml:"access_life_time"`
		Secret          string `yaml:"secret"`
	} `yaml:"auth"`
	Api struct {
		Url        string `yaml:"url"`
		ApiPublic  string `yaml:"api_public"`
		ApiPrivate string `yaml:"api_private"`
	}
}

func LoadConfig(path string) (*Config, error) {
	// #nosec
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file not found %s %w", path, err)
	}
	yamlDecoder := yaml.NewDecoder(file)

	cfg := &Config{}
	err = yamlDecoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("%v %w", err.Error())
	}

	return cfg, nil
}
