package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load(".env")
	if err != nil {
		return errors.Wrap(err, "environment variables")
	}

	db, err := setupDatabase()
	if err != nil {
		return errors.Wrap(err, "setup database")
	}

	srv := newServer()
	srv.db = db

	return http.ListenAndServe(":8080", srv)
}

// The server type contains the dependencies of our server.
type server struct {
	router *mux.Router
	db     *sql.DB
}

// newServer instantiates a server type and sets up its routes.
// Dependencies are not set up here so that it is easier to test.
func newServer() *server {
	srv := &server{
		router: mux.NewRouter(),
	}
	srv.routes()
	return srv
}

// Implementing ServeHTTP turns the server type into a http.Handler.
// Hence, server can be used wherever http.Handler can (e.g. http.ListenAndServe).
// Inside, we simply pass the execution to the router.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
