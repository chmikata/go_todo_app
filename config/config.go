package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Env        string `env:"TODO_ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"TODO_DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"TODO_DB_PORT" envDefault:"5432"`
	DBUser     string `env:"TODO_DB_USER" envDefault:"todo"`
	DBPassword string `env:"TODO_DB_PASSWORD" envDefault:"todo"`
	DBName     string `env:"TODO_DB_NAME" envDefault:"todo"`
	RedisHost  string `env:"TODO_REDIS_HOST" envDefault:"localhost"`
	RedisPort  int    `env:"TODO_REDIS_PORT" envDefault:"36379"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
