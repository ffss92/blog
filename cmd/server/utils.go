package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
)

// Walks the list of roots and collects all subdirectories.
func collectDirs(roots ...string) ([]string, error) {
	dirs := make([]string, 0)
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return dirs, nil
}

// Cache files with .ttf extension, 'no-cache' all others in development mode.
func cacheControl(r *http.Request) string {
	// Always cache fonts
	if filepath.Ext(r.URL.Path) == ".ttf" {
		return "public, max-age=31536000, immutable"
	}
	return "no-cache"
}
