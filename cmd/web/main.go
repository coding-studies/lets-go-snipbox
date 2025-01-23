package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Take command line argument `-addr`. E.g.:
	//
	//   $ go run ./cmd/web -addr=":4004"
	//   $ go run ./cmd/web -addr=":80"
	//
	// Uses ":4000" that cmdline option -addr is not provided.
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on " + *addr)

	err := http.ListenAndServe(*addr, mux)

	log.Fatal(err)
}
