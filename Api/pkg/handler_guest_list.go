package pkg

import (
	"encoding/json"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"io/ioutil"
	"net/http"
)

func (app *App) addGuestListsRoutes(routes conf.RoutesConfig){
	app.Router.Mux.HandleFunc(routes.GetGuestListUri, app.handleGuestListGet()).Methods("GET")
	app.Router.Mux.HandleFunc(routes.PostGuestListUri, app.handleGuestListPost()).Methods("POST")
	app.Router.Mux.HandleFunc(routes.DeleteGuestListUri, app.handleGuestListDelete()).Methods("DELETE")
}

func (app *App) handleGuestListGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

type Guest struct {
	Name               string `json:"name"`
	Table              int    `json:"table"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

func (app *App) handleGuestListPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		guest := Guest{}
		_ = json.Unmarshal(body, &guest)

		println(guest.Name, guest.Table, guest.AccompanyingGuests)
	}
}

func (app *App) handleGuestListDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
