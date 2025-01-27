package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	var (
		pgHost    = ""
		pgUser    = ""
		pgPass    = ""
		pgDBName  = ""
		pgSSLMode = "disable"
	)

	pgHost = os.Getenv("PG_HOST")
	pgUser = os.Getenv("PG_USER")
	pgPass = os.Getenv("PG_PASS")
	pgDBName = os.Getenv("PG_DBNAME")

	if os.Getenv("PG_SSLMODE") != "" {
		pgSSLMode = os.Getenv("PG_SSLMODE")
	}

	dsn := flag.String(
		"dsn",
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=%s",
			pgHost,
			pgUser,
			pgPass,
			pgDBName,
			pgSSLMode,
		),
		"Postgres Data Source Name",
	)

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

	db, err := openDB(*dsn)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("port", *addr))

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())

	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
