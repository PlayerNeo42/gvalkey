package config

import (
	"reflect"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Host     string `env:"GVK_HOST" envDefault:"0.0.0.0" validate:"required,hostname|ip"`
	Port     int    `env:"GVK_PORT" envDefault:"6379" validate:"required,min=1,max=65535"`
	LogLevel string `env:"GVK_LOG_LEVEL" envDefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
}

func Load() (*Config, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	c.LogLevel = strings.ToUpper(c.LogLevel)

	if err := validateConfig(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func validateConfig(c *Config) error {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("env")
		if name == "" {
			return fld.Name
		}
		return name
	})

	if err := v.Struct(c); err != nil {
		return err
	}
	return nil
}
