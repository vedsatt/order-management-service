package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"

	redis "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository/cache"
	postgres "gitlab.crja72.ru/golang/2025/spring/course/students/268295-aisavelev-edu.hse.ru-course-1478/internal/repository/database"
)

type Config struct {
	GrpcPort    string `env:"GRPC_PORT"    env-default:"50051"`
	GatewayPort string `env:"GATEWAY_PORT" env-default:"8080"`
	Environment string `env:"ENV"          env-default:"prod"`

	postgres.PostgresCfg
	redis.RedisCfg
}

func New() (*Config, error) {
	var cfg Config
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./config/.env"
	}

	err := cleanenv.ReadConfig(envPath, &cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
