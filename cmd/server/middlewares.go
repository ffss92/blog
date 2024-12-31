package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.cfg.dev {
			fmt.Println(r.Header.Get("X-Forwarded-For"))
		}
		next.ServeHTTP(w, r)
	})
}
