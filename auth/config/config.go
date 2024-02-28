package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	Service struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Version string `yaml:"version"`
		Name    string `yaml:"name"`
	} `yaml:"service"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
		SSLMode  string `yaml:"ssl_mode"`
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
	SecretsPath map[string]struct {
		ApiPublic  string `yaml:"api_public"`
		ApiPrivate string `yaml:"api_private"`
	} `yaml:"secrets_path"`
	Secrets map[string]struct {
		ApiPublic  string `yaml:"-"`
		ApiPrivate string `yaml:"-"`
	} `yaml:"-"`
	Core struct {
		Url  string `yaml:"url"`
		Name string `yaml:"name"`
	} `yaml:"core"`
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
		return nil, fmt.Errorf("%v", err.Error())
	}

	err = loadSecrets(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadSecrets(cfg *Config) error {
	for k, v := range cfg.SecretsPath {
		apiPublic, err := readFile(v.ApiPublic)
		if err != nil {
			return err
		}
		apiPrivate, err := readFile(v.ApiPrivate)
		if err != nil {
			return err
		}

		cfg.Secrets[k] = struct {
			ApiPublic  string `yaml:"-"`
			ApiPrivate string `yaml:"-"`
		}(struct {
			ApiPublic  string
			ApiPrivate string
		}{
			ApiPublic:  apiPublic,
			ApiPrivate: apiPrivate,
		})
	}

	return nil
}

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = f.Close()
	}()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
