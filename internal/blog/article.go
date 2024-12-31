package blog

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"slices"
	"strings"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

type ArticleMetadata struct {
	Title    string   `yaml:"title"`
	Subtitle string   `yaml:"subtitle"`
	Author   string   `yaml:"author"`
	Draft    bool     `yaml:"draft"`
	Date     string   `yaml:"date"`
	Tags     []string `yaml:"tags"`
}

type Article struct {
	Slug       string
	Content    template.HTML
	RawContent []byte
	PageViews  int

	ArticleMetadata
}

func (s *Service) GetArticle(ctx context.Context, slug string) (*Article, error) {
	article, err := s.getArticle(slug)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (s *Service) getArticle(slug string) (*Article, error) {
	if s.dev {
		if err := s.refreshArticles(); err != nil {
			return nil, err
		}
	}

	article, ok := s.cache[slug]
	if !ok {
		return nil, ErrArticleNotFound
	}
	if !s.dev && article.Draft {
		return nil, ErrArticleNotFound
	}
	return article, nil
}

func (s *Service) ListArticles(ctx context.Context) ([]*Article, error) {
	articles, err := s.listArticles()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *Service) listArticles() ([]*Article, error) {
	if s.dev {
		if err := s.refreshArticles(); err != nil {
			return nil, err
		}
	}

	articles := make([]*Article, 0)
	for _, article := range s.cache {
		if !s.dev && article.Draft {
			continue
		}
		articles = append(articles, article)
	}

	slices.SortFunc(articles, func(a *Article, b *Article) int {
		return -strings.Compare(a.Date, b.Date)
	})

	return articles, nil
}

func (s *Service) SavePageview(ctx context.Context, slug, ipAddress, userAgent, referrer string) error {
	query := `
	INSERT INTO pageviews (slug, ip_address, user_agent, referrer)
	VALUES ($1, $2, $3, $4)`
	args := []any{slug, ipAddress, userAgent, referrer}

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("blog: failed to save pageview: %w", err)
	}

	return nil
}
