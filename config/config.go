package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type Environment string

func (e Environment) IsGoogle() bool {
	return e == "google"
}

type Config struct {
	Env      Environment `env:"ENV"`
	Telegram TelegramConfig
	Google   GoogleConfig
	Logger   LoggerConfig
}

type TelegramConfig struct {
	Token   string `env:"TELEGRAM_TOKEN"`
	Timeout uint   `env:"TELEGRAM_POLLER_TIMEOUT, default=10"`
}

type GoogleConfig struct {
	ProjectID string `env:"GOOGLE_PROJECT_ID"`
	Secret    string `env:"GOOGLE_APPLICATION_CREDENTIALS"`
}

type LoggerConfig struct {
	Level string `env:"LOGGER_LEVEL,default=info"`
}

func (c *Config) validate() error {
	errT := c.validateTelegram()
	errG := c.validateGoogle()

	return multierr.Combine(errT, errG)
}

func (c *Config) validateTelegram() error {
	if c.Telegram.Token == "" {
		return errors.New("TELEGRAM_TOKEN: required")
	}

	return nil
}

func (c *Config) validateGoogle() error {
	if c.Google.ProjectID == "" {
		return errors.New("GOOGLE_PROJECT_ID: required")
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
