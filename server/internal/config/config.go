package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
)

type Config struct {
	AppVersion string `yaml:"app-version" env:"APP_VERSION" env-required:"true"`
	Env        string `yaml:"env" env:"ENV" env-default:"production"`
	DbUri      string `yaml:"db-uri" env:"DB_URI" env-required:"true"`
}

var cfg *Config

func MustLoad() *Config {
	if cfg != nil {
		return cfg
	}

	cfg = &Config{}

	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			slog.Error("config path is not set")
			os.Exit(1)
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("config file does not exist: %s", configPath)
		os.Exit(1)
	}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		slog.Error("cannot read config file: %s", err.Error())
		os.Exit(1)
	}

	return cfg
}
