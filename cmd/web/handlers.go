package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// home is a handler which writes a byte slicing simple text as the
// response body.
func (a *application) home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Server", "Go")

	templates := []string{
		"./cmd/web/ui/html/base.tmpl.html",
		"./cmd/web/ui/html/partials/nav.tmpl.html",
		"./cmd/web/ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// snippetView is a handler for viewing a specific snippet.
func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract and validate the id wildcard path parameter.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d", id)
}

// snippetCreate displays a form for creating a snippet.
func (a *application) snippetCreate(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}

// snippetCreatePost processes the post request to create a new snippet.
func (a *application) snippetCreatePost(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet"))
}
