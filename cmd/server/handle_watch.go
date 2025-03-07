package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (app *application) handleWatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targets, err := collectDirs("web/views", "web/static")
		if err != nil {
			app.serverError(w, r, err)
			return
		}

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
		err = rc.SetWriteDeadline(time.Time{})
		if err != nil {
			app.serverError(w, r, err)
			return
		}

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
				case fsnotify.Write, fsnotify.Create:
					fmt.Fprint(w, "event: mod\n")
					fmt.Fprintf(w, "data: %s\n\n", msg.Name)
					if err := rc.Flush(); err != nil {
						app.clientError(w, r, fmt.Errorf("failed to flush: %w", err))
						return
					}
				}
			}
		}
	}
}
