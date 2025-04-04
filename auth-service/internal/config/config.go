package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	StorageAddress string        `yaml:"storage_address"`
	TockenTTL      time.Duration `yaml:"tocken_ttl"`
	GRPC           GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	var cfg Config
	path := getFilePath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file does not exist on designated path")
	}

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("Failed to read config" + err.Error())
	}

	return &cfg
}

func getFilePath() string {
	var filePath string
	flag.StringVar(&filePath, "--logsFilePath", "", "path to config file")
	flag.Parse()

	if filePath == "" {
		filePath = os.Getenv("CONFIG_FILE")
	}

	return filePath
}
