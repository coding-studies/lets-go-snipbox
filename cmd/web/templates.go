package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snipbox.fernandobasso.dev/internal/models"
)

// tmplData is used as a type-safe way to pass multiple pieces of data
// to templates (since templates can only take a single variable containing
// data). By using a struct, we can pass a "bag" of composite data as needed.
type tmplData struct {
	// Calculated dynamically so it will automatically reflect each year without
	// requiring changes to the code. It is supposed to allways reflect the
	// actual, current year.
	CurrentYear int

	// For displaying a single snippet.
	Snippet models.Snippet

	// For displaying a list of snippets.
	Snippets []models.Snippet

	// Contains field validation errors (if any) with the submitted data.
	Form any

	// Contains a flash message.
	Flash string

	IsAuthenticated bool

	CSRFToken string
}

// humanDateAndTime returns a nice, readable representation of time.Time.
func humanDateAndTime(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var viewHelpers = template.FuncMap{
	"humanDateAndTime": humanDateAndTime,
}

// newTmplCache returns an in-memory cache of template sets.
//
// To retrieve a template from the cache we must provide the basename of the
// page template we want to use, like "home.tmpl.html" or "view.tmpl.html".
func newTmplCache() (map[string]*template.Template, error) {
	// The key is the basename of the template file, for example "home.tmpl.html"
	// or "view.tmpl.html"
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./cmd/web/ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		baseName := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.New(page).
			Funcs(viewHelpers).
			ParseFiles("./cmd/web/ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on the current template set to add partials.
		ts, err = ts.ParseGlob("./cmd/web/ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() on the current template set to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the current template set to the cache.
		cache[baseName] = ts
	}

	return cache, nil
}
