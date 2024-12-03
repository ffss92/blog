package main

import (
	"errors"
	"net/http"

	"ffss.dev/internal/blog"
	"github.com/go-chi/chi/v5"
)

type authorPage struct {
	basePage
	Author *blog.Author
}

func (app *application) handleAuthorShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle := chi.URLParam(r, "handle")
		author, err := app.blog.GetAuthor(r.Context(), handle)
		if err != nil {
			switch {
			case errors.Is(err, blog.ErrAuthorNotFound):
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		app.render(w, r, "authors/show", authorPage{
			basePage: app.newBasePage(r, author.Name),
			Author:   author,
		})
	}
}
