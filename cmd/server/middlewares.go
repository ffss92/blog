package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func (app *application) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.cfg.dev {
			r.RemoteAddr = strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		defer func() {
			app.logger.Info(
				"http request",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.RequestURI()),
				slog.String("ip_address", r.RemoteAddr),
				slog.Duration("duration", time.Since(t)),
			)
		}()
		next.ServeHTTP(w, r)
	})
}
