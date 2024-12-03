package templates

import (
	"html/template"
	"io/fs"
	"path/filepath"
)

func Parse(views fs.FS) (map[string]*template.Template, error) {
	viewSet := [][]string{
		{"home"},
		{"articles/index", "articles"},
		{"articles/show", "articles"},
		{"authors/show", "authors"},
	}

	cache := make(map[string]*template.Template)
	for _, set := range viewSet {
		tmpl := template.New("base.tmpl").Option("missingkey=zero")
		tmpl, err := tmpl.ParseFS(views, "*.tmpl")
		if err != nil {
			return nil, err
		}
		for _, f := range set {
			_, err = tmpl.ParseFS(views, filepath.Join(f, "*.tmpl"))
			if err != nil {
				return nil, err
			}
		}
		cache[set[0]] = tmpl
	}
	return cache, nil
}
