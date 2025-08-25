package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexey-dobry/auth-service/internal/domain/jwt"
	pg "github.com/alexey-dobry/auth-service/internal/repository/postgresql"
	"github.com/alexey-dobry/auth-service/internal/server/grpc"
	"github.com/alexey-dobry/auth-service/pkg/logger/zap"
	"github.com/alexey-dobry/auth-service/pkg/validator"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger     zap.Config  `yaml:"logger"`
	GRPC       grpc.Config `yaml:"grpc"`
	Repository pg.Config   `yaml:"repository"`
	JWT        jwt.Config  `yaml:"jwt"`
}

func MustLoad() Config {
	var cfg Config
	configPath := ParseFlag(cfg)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		errMsg := fmt.Sprintf("Failed to read config on path(%s): %s", configPath, err)
		panic(errMsg)
	}

	if err := validator.V.Struct(&cfg); err != nil {
		errMsg := fmt.Sprintf("Failed to validate config: %s", err)
		panic(errMsg)
	}

	return cfg
}

func ParseFlag(cfg Config) string {
	configPath := flag.String("config", "./config/config.yaml", "config file path")
	configHelp := flag.Bool("help", false, "show configuration help")

	if *configHelp {
		headerText := "Configuration options:"
		help, err := cleanenv.GetDescription(&cfg, &headerText)
		if err != nil {
			errMsg := fmt.Sprintf("error getting configuration description: %s", err)
			panic(errMsg)
		}
		fmt.Println(help)
		os.Exit(0)
	}

	return *configPath
}
