package sentry

import (
	"bytes"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/spiral/roadrunner"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	rrhttp "github.com/spiral/roadrunner/service/http"
	"net/http"
)

const ID = "sentry"

type Service struct {
}

func (s *Service) Init(cfg *Config) (bool, error) {
	err := raven.SetDSN(cfg.DSN)
	if err != nil {
		return false, err
	}
	svc, _ := rr.Container.Get(rrhttp.ID)
	if svc, ok := svc.(*rrhttp.Service); ok {
		svc.AddListener(systemListener)
		svc.AddListener(httpListener)
	}
	return true, nil
}

// listener listens to http events and generates nice looking output.
func httpListener(event int, ctx interface{}) {
	// http events
	switch event {
	case rrhttp.EventError:
		e := ctx.(*rrhttp.ErrorEvent)

		buf := new(bytes.Buffer)
		buf.ReadFrom(e.Request.Body)
		body := buf.String()

		meta := map[string]string{
			"event":       "EventError",
			"code":        "500",
			"method":      e.Request.Method,
			"site":        e.Request.Host,
			"uri":         uri(e.Request),
			"requestBody": body,
		}
		raven.CaptureErrorAndWait(e.Error, meta)
		return
	}
}

func systemListener(event int, ctx interface{}) {
	switch event {
	case roadrunner.EventWorkerError:
		e := ctx.(roadrunner.WorkerError)

		meta := map[string]string{
			"event": "EventWorkerError",
		}
		raven.CaptureErrorAndWait(e.Caused, meta)
		return
	case roadrunner.EventWorkerDead:
		e := ctx.(roadrunner.WorkerError)

		meta := map[string]string{
			"event": "EventWorkerDead",
		}
		raven.CaptureErrorAndWait(e.Caused, meta)
		return
	case roadrunner.EventPoolError:
		e := ctx.(roadrunner.WorkerError)

		meta := map[string]string{
			"event": "EventPoolError",
		}
		raven.CaptureErrorAndWait(e.Caused, meta)
		return
	case roadrunner.EventWorkerKill:
		e := ctx.(roadrunner.WorkerError)

		meta := map[string]string{
			"event": "EventWorkerKill",
		}
		raven.CaptureErrorAndWait(e.Caused, meta)
		return
	}
}

func uri(r *http.Request) string {
	if r.TLS != nil {
		return fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
	}

	return fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
}
