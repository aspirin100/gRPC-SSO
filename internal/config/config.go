package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
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

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

// helpful for tests.
func MustLoadByPath(configPath string) *Config {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		panic("config file doesn't exist: " + configPath)
	}

	var config Config

	err = cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		panic("failed to read config " + err.Error())
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
