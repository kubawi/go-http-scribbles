package main

import (
	"fmt"
	"net/http"
	"os"
)

// ignore me
type myDB struct{}

type service struct {
	db     *myDB
	router *http.ServeMux
	// can add other things here, like a logger
}

// this makes `service` implement http.Handler interface
func (s service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *service) routes() {
	s.router.HandleFunc("/api/", s.handleAPI())
	s.router.HandleFunc("/", s.handleIndex())
}

func (s *service) handleAPI() http.HandlerFunc {
	thing := []byte("Pretend this is a result of a query")
	// you can make your DB calls here, etc.
	// you can also pass handler-specific dependencies as arguments to this handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// and we use thing here
		w.Write(thing)
	}
}

func (s *service) handleIndex() http.HandlerFunc {
	message := "Hello, world!"
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(message))
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, dbtidy, err := setupDatabase()
	if err != nil {
		return fmt.Errorf("error setting up the database: %w ", err)
	}
	defer dbtidy()

	service := service{
		db:     db,
		router: http.NewServeMux(),
	}
	service.routes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: service,
	}

	return server.ListenAndServe()
}

func setupDatabase() (*myDB, func(), error) {
	// return nil, func() { fmt.Println("Closed DB connection") }, errors.New("database blew up")
	return &myDB{}, func() { fmt.Println("Closed DB connection") }, nil
}
