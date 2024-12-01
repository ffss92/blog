package blog

import (
	"context"
	"database/sql"
)

type SearchResult struct {
	Slug    string `json:"slug"`
	Snippet string `json:"snippet"`
}

func (s *Service) Search(ctx context.Context, term string) {
	return
}

func (s *Service) indexContents() error {
	return indexContents(s.db, s.cache)
}

func indexContents(db *sql.DB, articles map[string]*Article) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM blog_posts_fts")
	if err != nil {
		return err
	}

	for _, article := range articles {
		args := []any{article.Slug, article.Title, article.Subtitle, article.Content}
		_, err = tx.Exec("INSERT INTO blog_posts_fts (slug, title, subtitle, content) VALUES ($1, $2, $3, $4)", args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
