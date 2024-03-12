package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	Auth        `yaml:"auth"`
	Listen      `yaml:"listen"`
}

type Auth struct {
	HashSalt   string        `yaml:"hash_salt"`
	SigningKey string        `yaml:"signing_key"`
	TokenTTL   time.Duration `yaml:"token_ttl"`
}

type Listen struct {
	Type        string        `yaml:"type" env-default:"tcp"`
	Host        string        `yaml:"host" env-default:"127.0.0.1"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}
	return cfg
}
