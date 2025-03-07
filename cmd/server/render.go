package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/mileusna/useragent"
)

type basePage struct {
	Dev       bool
	IsMac     bool
	HTMLTitle string
	UserAgent useragent.UserAgent
}

func (app *application) newBasePage(r *http.Request, title string) basePage {
	ua := useragent.Parse(r.UserAgent())

	return basePage{
		Dev:       app.isDev(),
		HTMLTitle: title,
		IsMac:     ua.IsMacOS() || ua.IsIOS(),
		UserAgent: useragent.Parse(r.UserAgent()),
	}
}

func (app *application) render(w http.ResponseWriter, _ *http.Request, templateName string, data any) {
	buf, err := app.executeTemplate(templateName, data)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, _ = buf.WriteTo(w)
}

func (app *application) executeTemplate(name string, data any) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	tmpl, err := app.findTemplate(name)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (app *application) findTemplate(name string) (*template.Template, error) {
	if app.isDev() || app.templates == nil {
		app.mu.Lock()
		defer app.mu.Unlock()

		err := app.parseTemplates()
		if err != nil {
			return nil, err
		}
	}

	tmpl, ok := app.templates[name]
	if !ok {
		return nil, fmt.Errorf("failed to find template %q", name)
	}
	return tmpl, nil
}

func (app *application) parseTemplates() error {
	viewSet := [][]string{
		{"error"},
		{"authors/show", "authors"},
		{"articles/index", "articles"},
		{"articles/show", "articles"},
	}

	app.templates = make(map[string]*template.Template)
	for _, set := range viewSet {
		tmpl := template.New("base.tmpl").Option("missingkey=zero")
		tmpl, err := tmpl.ParseFS(app.views, "*.tmpl")
		if err != nil {
			return err
		}
		for _, f := range set {
			_, err = tmpl.ParseFS(app.views, filepath.Join(f, "*.tmpl"))
			if err != nil {
				return err
			}
		}
		app.templates[set[0]] = tmpl
	}

	return nil
}
