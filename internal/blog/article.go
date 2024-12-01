package blog

import (
	"context"
	"errors"
	"html/template"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

type ArticleMetadata struct {
	Title  string   `yaml:"title"`
	Author string   `yaml:"author"`
	Draft  bool     `yaml:"draft"`
	Date   string   `yaml:"date"`
	Tags   []string `yaml:"tags"`
}

type Article struct {
	Slug       string
	Content    template.HTML
	RawContent []byte

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
	if err := s.refreshArticles(); err != nil {
		return nil, err
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
	if err := s.refreshArticles(); err != nil {
		return nil, err
	}
	articles := make([]*Article, 0)
	for _, article := range s.cache {
		if s.dev && article.Draft {
			continue
		}
		articles = append(articles, article)
	}
	return articles, nil
}
