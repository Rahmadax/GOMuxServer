package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type server struct {
	router *mux.Router
}

// Setup
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)

		os.Exit(1)
	}
}

func run() error {
	configFilePath := flag.String("config", "social/config/pre_live.yml", "Path to config file")
	config, err := conf.Get(*configFilePath)
	if err != nil {
		return err
	}

	server, err := newServer(config)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8080", server.router.mux)
}

func newServer(config *Configuration) (*server, error) {
	s := &server{}
	s.router = mux.NewRouter()

	s.routes(config)

	return s, nil
}


func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
