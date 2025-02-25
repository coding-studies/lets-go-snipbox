package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	"snipbox.fernandobasso.dev/internal/models"

	_ "github.com/lib/pq"
)

type application struct {
	logger      *slog.Logger
	snippets    *models.SnippetModel
	tmplCache   map[string]*template.Template
	formDecoder *form.Decoder
	sessMgr     *scs.SessionManager
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
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmplCache, err := newTmplCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessMgr := scs.New()
	sessMgr.Store = postgresstore.New(db)
	sessMgr.Lifetime = 12 * time.Hour

	app := &application{
		logger:      logger,
		snippets:    &models.SnippetModel{DB: db},
		tmplCache:   tmplCache,
		formDecoder: form.NewDecoder(),
		sessMgr:     sessMgr,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	logger.Info("starting server", slog.String("port", srv.Addr))

	err = srv.ListenAndServe()

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
