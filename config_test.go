package sentry

import (
	"encoding/json"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCfg struct{ cfg string }

func (cfg *mockCfg) Get(name string) service.Config  { return nil }
func (cfg *mockCfg) Unmarshal(out interface{}) error { return json.Unmarshal([]byte(cfg.cfg), out) }

func Test_Config_Hydrate(t *testing.T) {
	cfg := &mockCfg{`{"DSN": "https://u:p@example.com/sentry/1"}`}
	c := &Config{}

	assert.NoError(t, c.Hydrate(cfg))
}

func Test_Config_Hydrate_Error(t *testing.T) {
	cfg := &mockCfg{`{"enable": true,"DSN": https://u:p@example.com/sentry/1"}`}
	c := &Config{}

	assert.Error(t, c.Hydrate(cfg))
}

func Test_Config_Valid(t *testing.T) {
	assert.NoError(t, (&Config{DSN: "https://u:p@example.com/sentry/1"}).Valid())
	assert.Error(t, (&Config{DSN: "//u:p@example.com/sentry/1"}).Valid())
	assert.Error(t, (&Config{DSN: "test@example.com"}).Valid())
}
