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
		if !app.isDev() {
			r.RemoteAddr = strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
		}
		next.ServeHTTP(w, r)
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *logResponseWriter) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func (l *logResponseWriter) Unwrap() http.ResponseWriter {
	return l.ResponseWriter
}

func (app *application) reqLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		lw := &logResponseWriter{ResponseWriter: w, status: http.StatusOK}

		defer func() {
			app.logger.Info(
				"http request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int("status", lw.status),
				slog.String("proto", r.Proto),
				slog.String("ip_address", r.RemoteAddr),
				slog.Duration("duration", time.Since(t)),
			)
		}()

		next.ServeHTTP(lw, r)
	})
}
