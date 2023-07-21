package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS"`            //адрес и порт запуска сервиса
	DSN            string `env:"DATABASE_URI"`           //адрес подключения к базе данных
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"` //адрес системы расчёта начислений
	LogLevel       string `env:"LOG_LEVEL"`
}

const (
	serverAddressDefault = "localhost:8080"
	dsnDefault           = ""
	accrualAddressDefault
	logLevelDefault = "info"
)

func NewConfig() *Config {
	cfg := Config{}
	cfg.init()

	if err := env.Parse(&cfg); err != nil {
		//logging
		return nil
	}

	return &cfg
}

func (c *Config) init() {
	flag.StringVar(&c.ServerAddress, "a", serverAddressDefault, "Listen server address (default - :8080)")
	flag.StringVar(&c.DSN, "d", dsnDefault, "URI to database")
	flag.StringVar(&c.AccrualAddress, "r", accrualAddressDefault, "Accrual system address")
	flag.StringVar(&c.LogLevel, "l", logLevelDefault, "Log level")
	flag.Parse()
}
