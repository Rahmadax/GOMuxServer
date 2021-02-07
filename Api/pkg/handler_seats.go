package pkg

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"net/http"
)

func (app *App) addSeatsRoutes(routes conf.RoutesConfig){
	app.Router.Mux.HandleFunc(routes.CountEmptySeatsUri, app.handleCountEmptySeats()).Methods("GET")
}

type emptySeats struct {

}

func (app *App) handleCountEmptySeats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}