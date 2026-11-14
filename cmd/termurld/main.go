package main

import (
	sql2 "database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/JUSTINS88/termurl/internal/store"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port     string
	database dbConfig
}

type dbConfig struct {
	url         string
	maxIdleTime time.Duration
}

type Application struct {
	config config
	logger *slog.Logger
	db     *store.Queries
	cache  *LinkCache
}

func main() {
	Serve()
}

func Serve() {
	_ = godotenv.Load("./.env")

	// Database Config
	dbURL := os.Getenv("DB_URL")
	dbConn, err := sql2.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed Connect to Database:", err)
	}

	// Port Config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// App Config
	app := &Application{
		config: config{
			port: port,
			database: dbConfig{
				url:         dbURL,
				maxIdleTime: 5 * time.Minute,
			},
		},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		db:     store.New(dbConn),
		cache:  NewLinkCache(),
	}

	srv := &http.Server{
		Addr:         ":" + app.config.port,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
