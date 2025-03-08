package main

import (
	"net/http"

	"github.com/ffss92/fileserver"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewMux()

	r.Use(app.realIP)
	r.Use(app.reqLogger)
	r.Use(app.recoverer)

	r.NotFound(app.notFound)

	filesrv := fileserver.New(app.static, fileserver.WithCacheControlFunc(cacheControl))
	r.Mount("/static/", http.StripPrefix("/static/", filesrv))

	r.Group(func(r chi.Router) {
		r.Get("/", app.handleHome())
		r.Get("/articles", app.handleArticleIndex())
		r.Get("/articles/{slug}", app.handleArticleShow())
		r.Get("/authors/{handle}", app.handleAuthorShow())

		r.Route("/api", func(r chi.Router) {
			r.Get("/search", app.handleSearch())
		})

		// Static
		r.Get("/robots.txt", app.handleRobotsTxt())
		r.Get("/favicon.ico", app.handleFavicon())

		// Dev
		if app.isDev() {
			r.Get("/watch", app.handleWatch())
		}
	})

	return r
}
