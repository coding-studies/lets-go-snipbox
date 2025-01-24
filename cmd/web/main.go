package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

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

	////
	// Log as JSON and make logs include debugging logs as well instead of
	// the default behavior which is to only log from info and above. Also
	// includes the source location where the log is being issued from.
	//
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("port", *addr))

	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())

	os.Exit(1)
}
