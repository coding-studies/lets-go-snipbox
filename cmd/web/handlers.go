package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snipbox.fernandobasso.dev/internal/models"
	"snipbox.fernandobasso.dev/internal/validator"
)

// home is a handler which writes a byte slicing simple text as the
// response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTmplData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
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

	data := app.newTmplData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

// snippetCreate displays a form for creating a snippet.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTmplData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int

	// Embeds Validator so snippetCreateForm “inherits” all fields from it.
	validator.Validator
}

// snippetCreatePost processes the post request to create a new snippet.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Fails if body contains more than 8912 bytes. Error surfaces to ParseForm().
	r.Body = http.MaxBytesReader(w, r.Body, 8912)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	////
	// DESIGN: In the book, not validation on the size of the content is being
	// performed, even though there is a default limit of 10MB (and our own custom
	// limit of 8912 bytes) on the size of the request.
	//
	form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This value must be one of 1, 7 or 365")

	// Dump any errors in the response as plain text for now.
	if !form.Valid() {
		data := app.newTmplData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
