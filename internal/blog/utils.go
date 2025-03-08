package blog

import (
	"io/fs"
	"path/filepath"
)

// Collects all markdown (.md, .markdown) file paths from a [io/fs.FS].
func collectMarkdown(articles fs.FS) ([]string, error) {
	paths := make([]string, 0)
	err := fs.WalkDir(articles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if ext := filepath.Ext(path); ext == ".md" || ext == ".markdown" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
