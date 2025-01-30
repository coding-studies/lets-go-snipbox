package main

import (
	"html/template"
	"path/filepath"
	"snipbox.fernandobasso.dev/internal/models"
)

// templateData is used as a type-safe way to pass multiple pieces of data
// to templates (since templates can only take a single variable containing
// data). By using a struct, we can pass a "bag" of composite data as needed.
type templateData struct {
	// For displaying a single snippet.
	Snippet models.Snippet

	// For displaying a list of snippets.
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./cmd/web/ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		baseName := filepath.Base(page)

		files := []string{
			"./cmd/web/ui/html/base.tmpl.html",
			"./cmd/web/ui/html/partials/nav.tmpl.html",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[baseName] = ts
	}

	return cache, nil
}
