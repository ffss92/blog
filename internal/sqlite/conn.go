package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "modernc.org/sqlite"
)

func Connect(ctx context.Context, dbPath string) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	q := make(url.Values)
	q.Set("_pragma", "busy_timeout(5000)")
	q.Set("_pragma", "journal_mode(WAL)")
	q.Set("_pragma", "foreign_keys(ON)")
	q.Set("_pragma", "synchronous(NORMAL)")

	db, err := sql.Open("sqlite", fmt.Sprintf("%s?%s", dbPath, q.Encode()))
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
