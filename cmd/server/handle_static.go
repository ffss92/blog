package main

import "net/http"

func (app *application) handleRobotsTxt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/robots.txt")
	}
}

func (app *application) handleFavicon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/images/favicon.ico")
	}
}
