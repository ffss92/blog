package main

import (
	"io/fs"
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
