---
title: Easy Auto-Reload with Server-Sent Events in Go.
subtitle: |
  Reloading pages every time we make a change in development can be a little tiresome.
  Discover how SSE can streamline auto-reloading for a smoother development process in Go.
author: "@ffss"
draft: false
date: "2024-12-30"
tags:
  - Go
  - HTTP
  - SSE
---

## Quick intro to Server-Sent Events

Server-Sent Events (SSE) enables servers to push real-time updates to
clients over a simple HTTP connection. Unlike WebSockets, which allow for two-way communication,
SSE is a **one-way** communication method where the server sends updates, and the client listens.
This makes SSE an excellent choice for scenarios like live notifications, activity feeds, or auto-reloading pages.

[MDN](https://developer.mozilla.org) provides great documentation on the
[event stream format](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format).
If you are not familiar with it, you should check it out. Implementing a SSE endpoint in Go is very straighforward:

```go
func handleSSE(w http.ResponseWriter, r *http.Request) {
  flusher, ok := w.(http.Flusher)
  if !ok {
    http.Error(w, "sse not supported", http.StatusBadRequest)
    return
  }

  // Send headers to client
  w.Header().Set("Content-Type", "text/event-stream")
  w.Header().Set("Cache-Control", "no-store")
  w.Header().Set("Connection", "keep-alive")
  flusher.Flush() 

  ticker := time.NewTicker(time.Second)
  defer ticker.Stop()

  for {
    select {
    case <-r.Context().Done():
      // Connection closed the connection
      return
    case t := <-ticker.C:
      // Send data to client every tick
      fmt.Fprint("event: tick\n")
      fmt.Fprintf("data: %s\n\n", t)
      flusher.Flush() 
    }
  }
}
```

In the example above, we sent the current time every second to the client.
Note that we set `Content-Type: text/event-stream` so the browser knows this
is a SSE endpoint.

The client implementation is also simple, since all major browsers support SSE
through the [EventSource API](https://developer.mozilla.org/en-US/docs/Web/API/EventSource).

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>SSE</title>
  </head>
  <body>
    <div id="time"></div>
    <script>
      const es = new EventSource("/sse");
      const timeEl = document.getElementById("time");
      es.addEventListener("tick", (e) => {
        // Updates the text in div#time every second.
        timeEl.innerText = e.data; 
      });
    </script>
  </body>
</html>
```

Since we defined a event type in our HTTP handler using `event: tick`. We can
use `addEventListener` to target that specific event in our client.

## Implementing a file watcher endpoint

Now that we know how SSE works and how to implement it in Go, we can use libraries like
[fsnotify](https://github.com/fsnotify/fsnotify) to send notifications to the browser whenever
a file is modified.

### Watching file changes

First, let's add the `fsnotify` package to the project by running:

```bash
go get github.com/fsnotify/fsnotify
```

Now, let's create a new HTTP handler that listens the files changes and send
events to the client:

```go
func handleWatch(w http.ResponseWriter, r *http.Request) {
  // Dirs that will be watched
  targets := []string{"articles", "templates"} 

  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  defer watcher.Close()

  for _, target := range targets {
    if err := watcher.Add(target); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }

  flusher, ok := w.(http.Flusher)
  if !ok {
    http.Error(w, "sse not supported", http.StatusBadRequest)
    return
  }

  w.Header().Set("Content-Type", "text/event-stream")
  w.Header().Set("Cache-Control", "no-store")
  w.Header().Set("Connection", "keep-alive")
  flusher.Flush()

  for {
    select {
      case <-r.Context().Done():
        return
      case msg := <-watcher.Events:
        switch msg.Op {
        case fsnotify.Write:
          fmt.Fprint(w, "event: mod\n")
          fmt.Fprint(w, "data: reload\n\n")
          flusher.Flush()
        }
    }
  }
}
```

### Reloading pages

The client implementation is also very simple. We just reload the page every time we get a
new `mod` event:

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>SSE</title>
  </head>
  <body>
    <!-- Page content here -->
    <script>
      const es = new EventSource("/sse");
      es.addEventListener("mod", () => {
        location.reload();
      });
    </script>
  </body>
</html>
```

If you are using tailwindcss or some other frontend tool that needs some time 
to compile, you can add a `setTimeout` to delay the reload a little bit:

```html
<script>
  const es = new EventSource("/sse");
  let id;
  es.addEventListener("mod", () => {
    clearTimeout(id);
    setTimeout(() => {
      location.reload();
    }, 500);
  })
</script>
```

You can check out the implementation that I use for the development of this
blog [here](https://github.com/ffss92/blog/blob/main/cmd/server/handle_watch.go).

Thanks for reading!
