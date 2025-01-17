package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// home is a handler which writes a byte slicing simple text as the
// response body.
func home(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello Snipbox"))
}

// snippetView is a handler for viewing a specific snippet.
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract and validate the id wildcard path parameter.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific snippet with ID %d", id)
	w.Write([]byte(msg))
}

// snippetCreate displays a form for creating a snippet.
func snippetCreate(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view/{id}", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Print a log message to say that the server is starting.
	log.Print("starting server on :4000")

	// Use the http.ListenAndServe() function to start a new web server.
	// We pass in two parameters: the TCP network address to listen on (in
	// this case ":4000") and the servemux we just created. If
	// http.ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit. Note that any error
	// returned by http.ListenAndServe() is always non-nil.
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
