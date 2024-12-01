package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"ffss.dev/cmd/server/templates"
)

type basePage struct {
	Dev       bool
	HTMLTitle string
}

func (app *application) newBasePage(_ *http.Request, title string) basePage {
	return basePage{
		Dev:       app.cfg.dev,
		HTMLTitle: title,
	}
}

func (app *application) findTemplate(name string) (*template.Template, error) {
	if app.cfg.dev {
		app.mu.Lock()
		defer app.mu.Unlock()
		templates, err := templates.Parse(app.views)
		if err != nil {
			return nil, err
		}
		app.templates = templates
	}

	tmpl, ok := app.templates[name]
	if !ok {
		return nil, fmt.Errorf("failed to find template %q", name)
	}
	return tmpl, nil
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

func (app *application) render(w http.ResponseWriter, _ *http.Request, templateName string, data any) {
	buf, err := app.executeTemplate(templateName, data)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = buf.WriteTo(w)
}
