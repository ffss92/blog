package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) handleSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		res, err := app.blog.Search(r.Context(), q)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
