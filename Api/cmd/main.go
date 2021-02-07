package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/guest_list"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type server struct {
	router   *mux.Router
	dbClient *sql.DB
}

// Setup
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)

		os.Exit(1)
	}
}

func run() error {
	configFilePath := flag.String("config", "Api/config/pre_live.yml", "Path to config file")
	config, err := conf.GetConfig(*configFilePath)
	if err != nil {
		return err
	}

	server, err := newServer()
	if err != nil {
		return err
	}

	server.initialiseRoutes(config)

	println("Listening on port - " + config.Server.HostPort)
	return http.ListenAndServe(":"+config.Server.HostPort, server.router)
}

func newServer() (*server, error) {
	s := &server{}

	s.router = mux.NewRouter()


	dbClient, err := sql.Open("mySql", "")
	if err != nil {
		return nil, err
	}

	s.dbClient = dbClient

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) initialiseRoutes(config *conf.Configuration) {
	guestListHandler := guest_list.NewGuestListHandler()

	s.router.HandleFunc(config.Routes.GetGuestListUrl, pkg.handleGuestListGet()).Methods("GET")
	//s.router.HandleFunc("", s.handleGuestListPost()).Methods("POST")
	//s.router.HandleFunc("", s.handleGuestListDelete()).Methods("DELETE")

	//s.router.HandleFunc("/invitation/{name}", s.handleInvitationGet()).Methods("GET")
	//
	//s.router.HandleFunc("/guests/{name}", s.handleGuestsGet()).Methods("GET")
	//s.router.HandleFunc("/guests/{name}", s.handleGuestsCreate()).Methods("POST")
}
