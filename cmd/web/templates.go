package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Danvs60/snippetbox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get slice of all file paths matching pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract just file name from full path (e.g. home.tmpl)
		name := filepath.Base(page)

		// 1. create empty template set
		// 2. register funcmap with Funcs 
		// 3. parse as normal
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// parse all partial templates
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add to in-memory cache
		cache[name] = ts
	}

	return cache, nil
}
