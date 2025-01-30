package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"snipbox.fernandobasso.dev/internal/models"
	"strconv"
)

// home is a handler which writes a byte slicing simple text as the
// response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templates := []string{
		"./cmd/web/ui/html/base.tmpl.html",
		"./cmd/web/ui/html/partials/nav.tmpl.html",
		"./cmd/web/ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippets: snippets,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

// snippetView is a handler for viewing a specific snippet.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract and validate the id wildcard path parameter.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	templates := []string{
		"./cmd/web/ui/html/base.tmpl.html",
		"./cmd/web/ui/html/partials/nav.tmpl.html",
		"./cmd/web/ui/html/pages/view.tmpl.html",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippet: snippet,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// snippetCreate displays a form for creating a snippet.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	insertedSnippetID, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", insertedSnippetID), http.StatusSeeOther)
}

// snippetCreatePost processes the post request to create a new snippet.
func (app *application) snippetCreatePost(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save app new snippet"))
}
