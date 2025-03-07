package main

import (
	"net/http"
)

func (app *application) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/articles", http.StatusFound)
	}
}
