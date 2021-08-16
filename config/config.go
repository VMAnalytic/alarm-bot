package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/pkg/errors"
)

type Config struct {
	Telegram TelegramConfig
}

type TelegramConfig struct {
	Token   string `env:"TELEGRAM_TOKEN"`
	Timeout uint   `env:"TELEGRAM_POLLER_TIMEOUT, default=10"`
}

type LoggerConfig struct {
	Level string `env:"LOGGER_LEVEL,default=info"`
}

func (c *Config) validate() error {
	err := c.validateTelegram()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Config) validateTelegram() error {
	if c.Telegram.Token == "" {
		return errors.New("TELEGRAM_TOKEN: required")
	}

	return nil
}

func NewConfig() (*Config, error) {
	var (
		config Config
		err    error
	)

	err = envdecode.Decode(&config)
	if err != nil {
		if err != envdecode.ErrNoTargetFieldsAreSet {
			return nil, errors.WithStack(err)
		}
	}

	err = config.validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}
