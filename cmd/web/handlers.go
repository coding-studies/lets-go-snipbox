package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snipbox.fernandobasso.dev/internal/models"
	"snipbox.fernandobasso.dev/internal/validator"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

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
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`

	// Embeds Validator so snippetCreateForm “inherits” all fields from it.
	// The "-" tells the decoder to ignore this field.
	validator.Validator `form:"-"`
}

// snippetCreatePost processes the post request to create a new snippet.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Fails if body contains more than 8912 bytes. Error surfaces to ParseForm().
	r.Body = http.MaxBytesReader(w, r.Body, 8912)

	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
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

	app.sessMgr.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTmplData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare a zero-valued instance of userSignupForm
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(form.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(form.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTmplData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if errors.Is(err, models.ErrDuplicateEmail) {
		form.AddFieldError("email", "Email address is already in use")

		data := app.newTmplData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
	}

	app.sessMgr.Put(r.Context(), "flash", "Signup successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusOK)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a form for logging in a user... ")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user... ")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user... ")
}