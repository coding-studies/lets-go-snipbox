package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
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

func (_ *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
