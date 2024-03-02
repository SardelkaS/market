package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Secrets struct {
	ApiPublic  string `yaml:"api_public"`
	ApiPrivate string `yaml:"api_private"`
}

type Config struct {
	Service struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Version string `yaml:"version"`
	} `yaml:"service"`
	Postgres struct {
		Host                string `yaml:"host"`
		Port                string `yaml:"port"`
		User                string `yaml:"user"`
		Password            string `yaml:"password"`
		DBName              string `yaml:"db_name"`
		SSLMode             string `yaml:"ssl_mode"`
		MaxConns            string `yaml:"max_conns"`
		MinConns            string `yaml:"min_conns"`
		MaxConnLifetime     string `yaml:"max_conn_lifetime"`
		MaxConnIdleTime     string `yaml:"max_conn_idle_time"`
		HealthCheckDuration string `yaml:"health_check_duration"`
	} `yaml:"postgres"`
	TgBot struct {
		Token  string `yaml:"token"`
		ChatId int64  `yaml:"chat_id"`
	} `yaml:"tg_bot"`
	SecretsPath map[string]Secrets `yaml:"secrets_path"`
	Secrets     map[string]Secrets `yaml:"-"`
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

	err = loadSecrets(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadSecrets(cfg *Config) error {
	cfg.Secrets = make(map[string]Secrets)

	for k, v := range cfg.SecretsPath {
		apiPublic, err := readFile(v.ApiPublic)
		if err != nil {
			return err
		}
		apiPrivate, err := readFile(v.ApiPrivate)
		if err != nil {
			return err
		}

		cfg.Secrets[k] = Secrets{
			ApiPublic:  apiPublic,
			ApiPrivate: apiPrivate,
		}
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
