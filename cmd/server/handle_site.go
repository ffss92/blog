package main

import (
	"net/http"
)

type homePage struct {
	basePage
}

func (app *application) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/articles", http.StatusFound)
	}
}
