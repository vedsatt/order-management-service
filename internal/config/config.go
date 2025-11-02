package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	GrpcPort    string `env:"GRPC_PORT"    env-default:"50051"`
	GatewayPort string `env:"GATEWAY_PORT" env-default:"8080"`
	Environment string `env:"ENV"          env-default:"prod"`
}

func New() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig("./config/.env", &cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
