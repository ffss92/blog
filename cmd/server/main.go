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

func (app *application) isDev() bool {
	return app.cfg.dev
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

	var (
		static   = os.DirFS(cfg.static)
		articles = os.DirFS(cfg.articles)
		views    = os.DirFS(cfg.views)
	)

	logger := logging.NewLogger(cfg.dev)
	slog.SetDefault(logger)

	db, err := sqlite.Connect(context.Background(), cfg.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	blog, err := blog.New(cfg.dev, db, articles)
	if err != nil {
		return err
	}

	app := &application{
		cfg:    cfg,
		logger: logger,
		db:     db,
		blog:   blog,
		static: static,
		views:  views,
	}
	return app.serve()
}
