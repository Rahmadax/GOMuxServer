package pkg

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"net/http"
)

func (app *App) addSeatsRoutes(routes conf.RoutesConfig){
	app.Router.Mux.HandleFunc(routes.GetEmptySeatsUri, app.handleGetEmptySeats()).Methods("GET")
}

func (app *App) handleGetEmptySeats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}