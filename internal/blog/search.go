package blog

import (
	"context"
)

type SearchResult struct {
	Articles []*ArticleSearchResult `json:"articles"`
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

type ArticleSearchResult struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`

	score float64
}

func (s *Service) searchArticles(ctx context.Context, q string) ([]*ArticleSearchResult, error) {
	query := `
	SELECT 
		slug, 
		title, 
		subtitle,
		bm25(blog_posts_fts) AS score
	FROM blog_posts_fts
	WHERE blog_posts_fts MATCH $1 
	ORDER BY score ASC
	LIMIT 10`
	rows, err := s.db.QueryContext(ctx, query, q)
	if err != nil {
		return nil, err
	}

	articles := make([]*ArticleSearchResult, 0)
	for rows.Next() {
		var result ArticleSearchResult
		err := rows.Scan(&result.Slug, &result.Title, &result.Subtitle, &result.score)
		if err != nil {
			return nil, err
		}
		articles = append(articles, &result)
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
		args := []any{article.Slug, article.Title, article.Subtitle, article.Content}
		_, err = tx.Exec("INSERT INTO blog_posts_fts (slug, title, subtitle, content) VALUES ($1, $2, $3, $4)", args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
