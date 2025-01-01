---
title: Easy Auto-Reload with Server-Sent Events in Go
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

[MDN](https://developer.mozilla.org) provides a great explantion of the
[event stream format](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format).
If you are not familiar with it, give it a read!

## Example: Sending Ticks to the Client

### Server Implementation

Implementing a SSE endpoint in Go is very straighforward. In the example below,
we will send an event to the client every second using a `time.Ticker`:

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
      // Client closed the connection
      return
    case t := <-ticker.C:
      // Send data to client every tick
      fmt.Fprint(w, "event: tick\n")
      fmt.Fprintf(w, "data: %s\n\n", t)
      flusher.Flush() 
    }
  }
}
```

Rember to set the `Content-Type` header to `text/event-stream` to ensure that
the browser recognizes this as an SSE endpoint.

### Client Implementation

Using SSE on the client side is simple, thanks to the [EventSource API](https://developer.mozilla.org/en-US/docs/Web/API/EventSource),
which is supported by all major browsers.

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

Since we specified the event type `tick` in our Go handler, we can listen for it
using `addEventListener` on the client side.

## Implementing a File Watcher Endpoint

Now that we know how SSE works and how to implement it in Go, we can use libraries like
[fsnotify](https://github.com/fsnotify/fsnotify) to send notifications to the browser whenever
a file is modified.

### Watching File Changes

First, let's add the `fsnotify` package to the project by running the following command:

```bash
go get github.com/fsnotify/fsnotify
```

Now let's implement a handler that sends events to the client whenever files in specified directories are modified:

```go
func handleWatch(w http.ResponseWriter, r *http.Request) {
  // Dirs that will be watched,
  // adapt this to your project structure
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

### Reloading Pages

The client implementation is even simpler compared to the ticker example. We just reload the page every time we get a
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

If you use tools like TailwindCSS or others that require compilation time, add a delay using `setTimeout`
to make sure the page is reloaded after the compilation is done:

```html
<script>
  const es = new EventSource("/sse");

  let id;
  es.addEventListener("mod", () => {
    clearTimeout(id);
    setTimeout(() => {
      location.reload();
    }, 500);
  });
</script>
```

## Wrapping Up

With SSE and `fsnotify`, you can create a lightweight and effective auto-reloading mechanism for development.

Check out the implementation I use for my blog's development
[here](https://github.com/ffss92/blog/blob/main/cmd/server/handle_watch.go).

Thanks for reading!
