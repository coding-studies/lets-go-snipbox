package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	////
	// Take command line argument `-addr`. E.g.:
	//
	//   $ go run ./cmd/web -addr=":4004"
	//   $ go run ./cmd/web -addr=":80"
	//
	// Uses ":4000" that cmdline option -addr is not provided.
	//
	// And use this to get a list of command line flags this package supports:
	//
	//   $ go run ./cmd/web -help
	//
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on " + *addr)
	logger.Info("starting server", slog.String("port", *addr))

	err := http.ListenAndServe(*addr, mux)

	logger.Error(err.Error())

	os.Exit(1)
}
