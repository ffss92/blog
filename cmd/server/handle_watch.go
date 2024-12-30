package main

import (
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"
)

func (app *application) handleWatch() http.HandlerFunc {
	targets := []string{
		"articles",
		"web/views",
		"web/views/home",
		"web/views/authors",
		"web/views/authors/show",
		"web/views/articles",
		"web/views/articles/index",
		"web/views/articles/show",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		defer watcher.Close()

		for _, target := range targets {
			err := watcher.Add(target)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
		}

		rc := http.NewResponseController(w)
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Connection", "keep-alive")
		if err := rc.Flush(); err != nil {
			app.clientError(w, r, fmt.Errorf("failed to flush: %w", err))
			return
		}

		for {
			select {
			case <-r.Context().Done():
				return
			case msg := <-watcher.Events:
				switch msg.Op {
				case fsnotify.Write:
					fmt.Fprint(w, "event: mod\n")
					fmt.Fprint(w, "data: reload\n\n")
					if err := rc.Flush(); err != nil {
						app.clientError(w, r, fmt.Errorf("failed to flush: %w", err))
						return
					}
				}
			}
		}
	}
}
