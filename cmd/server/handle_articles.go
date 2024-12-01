package main

import (
	"errors"
	"net/http"

	"ffss.dev/internal/blog"
	"github.com/go-chi/chi/v5"
)

type articleIndexPage struct {
	basePage
	Articles []*blog.Article
}

func (app *application) handleArticleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles, err := app.blog.ListArticles(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		app.render(w, r, "articles/index", articleIndexPage{
			basePage: app.newBasePage(r, "Articles"),
			Articles: articles,
		})
	}
}

type articleShowPage struct {
	basePage
	Article *blog.Article
}

func (app *application) handleArticleShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		article, err := app.blog.GetArticle(r.Context(), slug)
		if err != nil {
			switch {
			case errors.Is(err, blog.ErrArticleNotFound):
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		app.render(w, r, "articles/show", articleShowPage{
			basePage: app.newBasePage(r, article.Title),
			Article:  article,
		})
	}
}
