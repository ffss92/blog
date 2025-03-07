package blog

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrAuthorNotFound = errors.New("author not found")
)

type Author struct {
	ID        int64  `json:"id"`
	Handle    string `json:"handle"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Birth     string `json:"birth"`
	ImageURL  string `json:"image_url"`
	GithubURL string `json:"github_url"`
}

// Grabs an author by it's handle. If handle is prefixed with and '@', it is trimmed automatically. Returns
// [ErrAuthorNotFound] if no results are returned.
func (b *Service) GetAuthor(ctx context.Context, handle string) (*Author, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	handle = strings.TrimPrefix(handle, "@")
	query := `
		SELECT 
			id, 
			handle, 
			name, 
			bio, 
			birth, 
			image_url, 
			github_url 
		FROM authors 
		WHERE handle = $1`

	var author Author
	err := b.db.QueryRowContext(ctx, query, handle).Scan(
		&author.ID,
		&author.Handle,
		&author.Name,
		&author.Bio,
		&author.Birth,
		&author.ImageURL,
		&author.GithubURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}

	return &author, nil
}
