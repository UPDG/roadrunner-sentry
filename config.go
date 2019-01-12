package sentry

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spiral/roadrunner/service"
)

type Config struct {
	// Contains Sentry DSN.
	DSN string
}

// Hydrate must populate Config values using given Config source. Must return error if Config is not valid.
func (c *Config) Hydrate(cfg service.Config) error {
	if err := cfg.Unmarshal(c); err != nil {
		return err
	}
	return c.Valid()
}

func (c *Config) Valid() error {
	if !govalidator.IsRequestURL(c.DSN) {
		return fmt.Errorf("sentry DSN in not URL")
	}
	return nil
}
