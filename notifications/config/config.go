package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string `yaml:"env" env-default:"local"`
	Kafka     `yaml:"kafka" env-required:"true"`
	WebSocket `yaml:"websocket"`
}

type Kafka struct {
	Brokers         []string `yaml:"brokers"`
	Topics          []string `yaml:"topics"`
	Address         string   `yaml:"address"`
	Port            int      `yaml:"port"`
	ConsumerGroupID string   `yaml:"group"`
}

type WebSocket struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
