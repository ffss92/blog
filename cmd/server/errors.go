package main

import (
	"log/slog"
	"net/http"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"unexpected error",
		slog.String("err", err.Error()),
		slog.String("method", r.Method),
		slog.String("uri", r.URL.RequestURI()),
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"bad request",
		slog.String("err", err.Error()),
		slog.String("uri", r.URL.RequestURI()),
	)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
