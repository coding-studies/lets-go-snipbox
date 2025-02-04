package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// serverError writes a log entry including useful request information and then
// sends a 500 to the client.
//
// DESIGN: Should we have a function that does X _and_ Y? In this case, log
// _and_ send HTTP error response?
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(
		err.Error(),
		slog.String("method", method),
		slog.String("uri", uri),
		slog.String("trace", trace),
	)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// newTmplData returns an instance of tmplData, initializing it
// with the current year.
func (app *application) newTmplData(_ *http.Request) tmplData {
	return tmplData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// render renders a given page.
//
// The page parameter refers to the basename of template we want to render, like
// "home.tmpl.html" or "edit.tmpl.html".
func (app *application) render(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	page string,
	data tmplData,
) {
	ts, ok := app.tmplCache[page]
	if !ok {
		err := fmt.Errorf("template '%s' does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	// Write the template to the buffer instead of directly to the response so
	// we can make sure there are no errors before attempting to write a
	// response to the user and thus presenting the user with the main page
	// skeleton filled with an error message. Either we present the user with
	// good, valid page view, or an error page, but not a mix.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If no error, then we can safely use the status code we get from
	// the parameter.
	w.WriteHeader(status)

	// And finally send the good page to the user.
	buf.WriteTo(w)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		// If we try to use an invalid/nil target destination, Decode() will return
		// an error of the type *formInvalidDecoderError.
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
	}

	// For all other errors, return them as usual.
	return err
}
