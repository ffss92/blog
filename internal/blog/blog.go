package blog

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strings"
	"sync"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

type Service struct {
	dev      bool
	db       *sql.DB
	md       goldmark.Markdown
	articles fs.FS

	mu    sync.Mutex
	cache map[string]*Article
}

func New(dev bool, db *sql.DB, articles fs.FS) (*Service, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("catppuccin-mocha"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(false),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	service := &Service{
		dev:      dev,
		db:       db,
		md:       md,
		articles: articles,
	}

	err := service.parseArticles()
	if err != nil {
		return nil, err
	}

	err = service.indexContents()
	if err != nil {
		return nil, err
	}

	return service, nil
}

// Collects all markdown files (.md) from a [fs.FS].
func markdownFiles(articles fs.FS) ([]string, error) {
	paths := make([]string, 0)
	err := fs.WalkDir(articles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(d.Name(), ".md") {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

func (s *Service) parseArticles() error {
	paths, err := markdownFiles(s.articles)
	if err != nil {
		return err
	}

	cache := make(map[string]*Article)
	for _, path := range paths {
		f, err := s.articles.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		buf := new(bytes.Buffer)
		contents, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		context := parser.NewContext()
		err = s.md.Convert(contents, buf, parser.WithContext(context))
		if err != nil {
			return err
		}

		metadata := meta.Get(context)
		b, err := yaml.Marshal(metadata)
		if err != nil {
			return err
		}

		var articleMetadata ArticleMetadata
		err = yaml.Unmarshal(b, &articleMetadata)
		if err != nil {
			return err
		}

		slug := strings.TrimSuffix(path, ".md")

		var views int
		err = s.db.QueryRow("SELECT COUNT(*) FROM pageviews WHERE slug = $1", slug).Scan(&views)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				views = 0
			default:
				return err
			}
		}

		cache[slug] = &Article{
			Slug:            slug,
			Content:         template.HTML(buf.String()),
			RawContent:      string(contents),
			ArticleMetadata: articleMetadata,
			PageViews:       views,
		}
	}

	s.cache = cache
	return nil
}

// If dev mode is on, parses all articles and set them to the cache.
func (s *Service) refreshArticles() error {
	if !s.dev {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.parseArticles()
	if err != nil {
		return fmt.Errorf("failed to refresh articles from fs: %w", err)
	}

	err = s.indexContents()
	if err != nil {
		return err
	}
	return nil
}
