package main

import (
	"errors"
	"fmt"
	"net/http"
	"snipbox.fernandobasso.dev/internal/models"
	"strconv"
)

// home is a handler which writes a byte slicing simple text as the
// response body.
func (a *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	//templates := []string{
	//	"./cmd/web/ui/html/base.tmpl.html",
	//	"./cmd/web/ui/html/partials/nav.tmpl.html",
	//	"./cmd/web/ui/html/pages/home.tmpl.html",
	//}
	//
	//ts, err := template.ParseFiles(templates...)
	//if err != nil {
	//	log.Print(err.Error())
	//	a.serverError(w, r, err)
	//	return
	//}
	//
	//err = ts.ExecuteTemplate(w, "base", nil)
	//if err != nil {
	//	a.serverError(w, r, err)
	//	return
	//}
}

// snippetView is a handler for viewing a specific snippet.
func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract and validate the id wildcard path parameter.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			a.serverError(w, r, err)
		}

		return
	}

	// Write snippet as plain/text for now.
	fmt.Fprintf(w, "%+v", snippet)
}

// snippetCreate displays a form for creating a snippet.
func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	insertedSnippetID, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", insertedSnippetID), http.StatusSeeOther)
}

// snippetCreatePost processes the post request to create a new snippet.
func (a *application) snippetCreatePost(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet"))
}
