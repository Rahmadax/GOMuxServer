package pkg

import (
	"net/http"
)

func (app *App) addGuestsRoutes(){
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.GetGuestsUri, app.handleGuestsGet()).Methods("GET")
	app.Router.HandleFunc(routes.PutGuestsUri, app.handleGuestsUpdate()).Methods("PUT")
}

func (app *App) handleGuestsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (app *App) handleGuestsUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
