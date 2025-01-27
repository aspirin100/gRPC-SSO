package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrEmptyPath      = errors.New("empty path to configuration file")
	ErrConfigNotFound = errors.New("config file not found")
)

type Config struct {
	Env         string        `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string        `yaml:"storagePath" env:"STORAGE_PATH" env-required:"true"`
	AccessTTL   time.Duration `yaml:"accessTokenTTL" env:"ACCESS_TTL" env-default:"60m"`      //nolint:tagliatelle
	RefreshTTL  time.Duration `yaml:"refreshTokenTTL" env:"REFRESH_TTL" env-default:"43200m"` //nolint:tagliatelle
	GRPC        GRPCConfig    `yaml:"grpc" env:"GRPC"`
	SecretKey   string        `env:"SECRET_KEY" env-required:"true"` // not safe to save in config file.
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env:"PORT" env-default:"8000"`
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"300m"`
}

func Load() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, ErrEmptyPath //nolint:wrapcheck
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, ErrConfigNotFound //nolint:wrapcheck
	}

	var config Config

	err = cleanenv.ReadConfig(path, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &config, nil
}

func MustLoadByPath(path string) *Config {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		panic(err)
	}

	var config Config

	err = cleanenv.ReadConfig(path, &config)
	if err != nil {
		panic(err)
	}

	return &config
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "confpath", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
