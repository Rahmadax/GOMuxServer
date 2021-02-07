package pkg

import (
	"net/http"
)

func (app *App) addInvitationRoutes(){
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.GetInvitationUri, app.handleInvitationGet()).Methods("GET")
}

func (app *App)handleInvitationGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
