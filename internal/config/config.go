package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerAddress     string        `env:"RUN_ADDRESS"`            //адрес и порт запуска сервиса
	DSN               string        `env:"DATABASE_URI"`           //адрес подключения к базе данных
	AccrualAddress    string        `env:"ACCRUAL_SYSTEM_ADDRESS"` //адрес системы расчёта начислений
	LogLevel          string        `env:"LOG_LEVEL"`
	SecretKey         string        `env:"JWT_SECRET"`
	RepositoryTimeout time.Duration `env:"REPOSITORY_TIMEOUT"`
	WorkersCount      int           `env:"WORKERS_COUNT"`
}

const (
	serverAddressDefault     = "localhost:8080"
	dsnDefault               = ""
	accrualAddressDefault    = "localhost:8090"
	logLevelDefault          = "info"
	secretKeyDefault         = "secret"
	repositoryTimeoutDefault = 5 * time.Second
	workersCountDefault      = 20
)

func NewConfig() (*Config, error) {
	cfg := Config{}
	cfg.init()

	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("error with parse config")
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) init() {
	flag.StringVar(&c.ServerAddress, "a", serverAddressDefault, "Listen server address (default - :8080)")
	flag.StringVar(&c.DSN, "d", dsnDefault, "URI to database")
	flag.StringVar(&c.AccrualAddress, "r", accrualAddressDefault, "Accrual system address")
	flag.StringVar(&c.LogLevel, "l", logLevelDefault, "Log level")
	flag.StringVar(&c.SecretKey, "s", secretKeyDefault, "Authorization token encryption key")
	flag.DurationVar(&c.RepositoryTimeout, "t", repositoryTimeoutDefault, "repository timeout")
	flag.IntVar(&c.WorkersCount, "w", workersCountDefault, "workers count")
	flag.Parse()
}
