package base

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type loggingHandler struct {
	handler http.Handler
	logger  zerolog.Logger
}

func (l loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	l.logger.Debug().Msgf("received '%s %s'", r.Method, path)
	elapsed := l.serveAndReportTime(w, r)
	l.logger.Debug().Msgf("request '%s %s' processed in %.3fms", r.Method, path, float64(elapsed.Nanoseconds())/1e6)
}

func (l loggingHandler) serveAndReportTime(w http.ResponseWriter, r *http.Request) time.Duration {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	return time.Since(start)
}
