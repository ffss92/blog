package main

import (
	"log/slog"
	"net/http"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:     app.cfg.addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting http server", slog.String("addr", app.cfg.addr))
	return srv.ListenAndServe()
}
