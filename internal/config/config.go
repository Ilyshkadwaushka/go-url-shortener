package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"prod" env-required:"true"` // env-default - значение по умолчанию
	StoragePath string `yaml:"storage_path" env-default:"./storage" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `env:"HTTP_PASSWORD" env-required:"true"`
}

func MustLoad() *Config { // Must - вместо возврата ошибки, функция паникует, а не возвращает ошибку
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file doesn`t exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
