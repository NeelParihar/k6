package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/api/common"
	v1 "go.k6.io/k6/api/v1"
	"go.k6.io/k6/core"
)

func newHandler(logger logrus.FieldLogger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/v1/", v1.NewHandler())
	mux.Handle("/ping", handlePing(logger))
	mux.Handle("/", handlePing(logger))
	return mux
}

// GetServer returns a http.Server instance that can serve k6's REST API.
func GetServer(addr string, engine *core.Engine, logger logrus.FieldLogger) *http.Server {
	mux := withEngine(engine, newLogger(logger, newHandler(logger)))
	return &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w wrappedResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// newLogger returns the middleware which logs response status for request.
func newLogger(l logrus.FieldLogger, next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		wrapped := wrappedResponseWriter{ResponseWriter: rw, status: 200} // The default status code is 200 if it's not set
		next.ServeHTTP(wrapped, r)

		l.WithField("status", wrapped.status).Debugf("%s %s", r.Method, r.URL.Path)
	}
}

func withEngine(engine *core.Engine, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		r = r.WithContext(common.WithEngine(r.Context(), engine))
		next.ServeHTTP(rw, r)
	})
}

func handlePing(logger logrus.FieldLogger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
		if _, err := fmt.Fprint(rw, "ok"); err != nil {
			logger.WithError(err).Error("Error while printing ok")
		}
	})
}
