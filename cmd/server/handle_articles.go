package main

import (
	"errors"
	"log/slog"
	"net/http"

	"ffss.dev/internal/blog"
	"github.com/go-chi/chi/v5"
)

type articleIndexPage struct {
	basePage
	Sort     string
	Articles []*blog.Article
}

func (app *application) handleArticleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sort string

		switch r.URL.Query().Get("sort") {
		case "date", "popular":
			sort = r.URL.Query().Get("sort")
		default:
			sort = "date"
		}

		articles, err := app.blog.ListArticles(r.Context(), sort)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.render(w, r, "articles/index", articleIndexPage{
			basePage: app.newBasePage(r, "Articles"),
			Articles: articles,
			Sort:     sort,
		})
	}
}

type articleShowPage struct {
	basePage
	Article *blog.Article
	Author  *blog.Author
}

func (app *application) handleArticleShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		article, err := app.blog.GetArticle(r.Context(), slug)
		if err != nil {
			switch {
			case errors.Is(err, blog.ErrArticleNotFound):
				app.notFound(w, r)
			default:
				app.serverError(w, r, err)
			}
			return
		}

		author, err := app.blog.GetAuthor(r.Context(), article.Author)
		if err != nil {
			switch {
			case errors.Is(err, blog.ErrAuthorNotFound):
				app.notFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = app.blog.SavePageview(r.Context(), slug, r.RemoteAddr, r.UserAgent(), r.Referer())
		if err != nil {
			app.logger.Error("failed to save pageview", slog.String("err", err.Error()))
		}

		app.render(w, r, "articles/show", articleShowPage{
			basePage: app.newBasePage(r, article.Title),
			Article:  article,
			Author:   author,
		})
	}
}
