package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env-default:"local"`
	Grpc gRPC   `yaml:"gRpc" env-required:"true"`
}

type gRPC struct {
	Port int `yaml:"port" env-required:"true"`
}

func MustReadCfg() *Config {
	path := fetchPathCfg()

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("%v", err)
	}

	cfg := Config{}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("%v", err)
	}

	return &cfg
}

var (
	pathCfg = flag.String("cfg", "./config/config.yaml", "config path")
)

func fetchPathCfg() string {
	flag.Parse()
	return *pathCfg
}
