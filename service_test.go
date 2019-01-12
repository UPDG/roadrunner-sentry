package sentry

import (
	"crypto/tls"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

type testCfg struct {
	sentryCgf string
	target    string
}

func (cfg *testCfg) Get(name string) service.Config {
	if name == ID {
		return &testCfg{target: cfg.sentryCgf}
	}

	return nil
}
func (cfg *testCfg) Unmarshal(out interface{}) error {
	return json.Unmarshal([]byte(cfg.target), out)
}

func Test_Service_Init_NoProject(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.Error(t, c.Init(&testCfg{sentryCgf: `{"DSN":"https://u:p@example.com/sentry/"}`}))
}

func Test_Service_Init_NoUser(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.Error(t, c.Init(&testCfg{sentryCgf: `{"DSN":"https://example.com/sentry"}`}))
}

func Test_Service_Init_AllGood(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{sentryCgf: `{"DSN":"https://u:p@example.com/sentry/1"}`}))
}

func TestService_Uri(t *testing.T) {
	testUrl, _ := url.Parse("/sentry/1")
	assert.Equal(t, "https://example.com/sentry/1", uri(&http.Request{
		TLS:  &tls.ConnectionState{},
		Host: "example.com",
		URL:  testUrl,
	}))
	assert.NotEqual(t, "http://example.com/sentry/1", uri(&http.Request{
		TLS:  &tls.ConnectionState{},
		Host: "example.com",
		URL:  testUrl,
	}))

	assert.Equal(t, "http://example.com/sentry/1", uri(&http.Request{
		Host: "example.com",
		URL:  testUrl,
	}))
	assert.NotEqual(t, "https://example.com/sentry/1", uri(&http.Request{
		Host: "example.com",
		URL:  testUrl,
	}))
}
