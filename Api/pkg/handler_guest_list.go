package pkg

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
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

func (app *App) handleGuestListPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (app *App) handleGuestListDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
