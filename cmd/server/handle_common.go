package main

import "net/http"

type homePage struct {
	basePage
}

func (app *application) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "home", homePage{
			basePage: app.newBasePage(r, "Home"),
		})
	}
}
