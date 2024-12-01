package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"sync"

	"ffss.dev/cmd/server/templates"
	"ffss.dev/internal/blog"
	"ffss.dev/internal/logging"
	"ffss.dev/internal/sqlite"
)

type config struct {
	addr     string
	dev      bool
	articles string
	static   string
	views    string
	dbPath   string
}

type application struct {
	cfg    config
	logger *slog.Logger
	db     *sql.DB
	blog   *blog.Service
	views  fs.FS
	static fs.FS

	mu        sync.Mutex
	templates map[string]*template.Template
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "Sets the HTTP server listen address.")
	flag.BoolVar(&cfg.dev, "dev", true, "Sets the application in development mode.")
	flag.StringVar(&cfg.articles, "articles", "articles", "Sets the articles dir.")
	flag.StringVar(&cfg.static, "static", "web/static", "Sets the static dir.")
	flag.StringVar(&cfg.views, "views", "web/views", "Sets the views dir.")
	flag.StringVar(&cfg.dbPath, "db-path", "blog.db", "Sets the sqlite database path.")
	flag.Parse()

	logger := logging.NewLogger(slog.LevelInfo, cfg.dev)

	db, err := sqlite.Connect(context.Background(), cfg.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	blog, err := blog.New(cfg.dev, db, os.DirFS(cfg.articles))
	if err != nil {
		return err
	}

	templates, err := templates.Parse(os.DirFS(cfg.views))
	if err != nil {
		return err
	}

	app := &application{
		cfg:    cfg,
		logger: logger,
		db:     db,
		blog:   blog,
		static: os.DirFS(cfg.static),
		views:  os.DirFS(cfg.views),

		mu:        sync.Mutex{},
		templates: templates,
	}
	return app.serve()
}
