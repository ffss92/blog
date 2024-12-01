package main

import (
	"net/http"

	"github.com/ffss92/fileserver"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewMux()

	r.Mount("/static/", http.StripPrefix("/static/", fileserver.ServeFS(app.static)))

	r.Group(func(r chi.Router) {
		r.Get("/", app.handleHome())
		r.Get("/articles", app.handleArticleIndex())
		r.Get("/articles/{slug}", app.handleArticleShow())
	})

	return r
}
