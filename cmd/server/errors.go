package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type errorPage struct {
	basePage
	StatusCode int
	StatusText string
	Message    string
}

func (app *application) renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	app.render(w, r, "error", errorPage{
		basePage:   app.newBasePage(r, fmt.Sprint(status)),
		StatusCode: status,
		StatusText: http.StatusText(status),
		Message:    message,
	})
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"unexpected error",
		slog.String("err", err.Error()),
		slog.String("method", r.Method),
		slog.String("uri", r.URL.RequestURI()),
	)
	app.renderError(w, r, http.StatusInternalServerError, "The server encountered an unexpected error serving your request.")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.renderError(w, r, http.StatusNotFound, "The requested page could not be found.")
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"bad request",
		slog.String("err", err.Error()),
		slog.String("uri", r.URL.RequestURI()),
	)
	app.renderError(w, r, http.StatusBadRequest, "Invalid or malformed request.")
}
