package pkg

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"net/http"
)

func (app *App) addGuestRoutes(routes conf.RoutesConfig){
	app.Router.Mux.HandleFunc(routes.GetGuestsUri, app.handleGuestsGet()).Methods("GET")
	app.Router.Mux.HandleFunc(routes.PutGuestsUri, app.handleGuestsUpdate()).Methods("PUT")
}

func (app *App) handleGuestsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (app *App) handleGuestsUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
