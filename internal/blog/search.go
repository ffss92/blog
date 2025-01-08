package blog

import (
	"context"
	"strconv"
)

type SearchResult struct {
	Articles []*ArticleResult `json:"articles"`
}

func (s *Service) Search(ctx context.Context, q string) (*SearchResult, error) {
	articles, err := s.searchArticles(ctx, q)
	if err != nil {
		return nil, err
	}

	result := &SearchResult{
		Articles: articles,
	}
	return result, nil
}

type ArticleResult struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`

	score float64
}

func (s *Service) searchArticles(ctx context.Context, q string) ([]*ArticleResult, error) {
	query := `
	SELECT 
		slug, 
		title, 
		subtitle,
		bm25(blog_posts_fts) AS score
	FROM blog_posts_fts
	WHERE blog_posts_fts MATCH $1 
	ORDER BY score ASC
	LIMIT 5`
	rows, err := s.db.QueryContext(ctx, query, strconv.Quote(q))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]*ArticleResult, 0)
	for rows.Next() {
		var result ArticleResult
		err := rows.Scan(&result.Slug, &result.Title, &result.Subtitle, &result.score)
		if err != nil {
			return nil, err
		}
		articles = append(articles, &result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return articles, nil
}

func (s *Service) indexContents() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM blog_posts_fts")
	if err != nil {
		return err
	}

	for _, article := range s.cache {
		// Skip draft articles in prod
		if article.Draft && !s.dev {
			continue
		}

		query := `
		INSERT INTO blog_posts_fts (slug, title, subtitle, content)
		VALUES ($1, $2, $3, $4)`
		args := []any{article.Slug, article.Title, article.Subtitle, article.Content}
		_, err = tx.Exec(query, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
