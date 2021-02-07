package pkg

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"net/http"
)

func (app *App) addInvitationRoutes(routes conf.RoutesConfig){
	app.Router.Mux.HandleFunc(routes.GetInvitationUri, app.handleInvitationGet()).Methods("GET")
}

func (app *App)handleInvitationGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
