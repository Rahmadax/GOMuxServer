package pkg

import (
"encoding/json"
"github.com/Rahmadax/GOMuxServer/Api/pkg/queries"
"net/http"
)

func (app *App) addGuestsRoutes() {
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.GetGuestsUri, app.getGuestsHandler()).Methods("GET")
	app.Router.HandleFunc(routes.PutGuestsUri, app.updateGuestsHandler()).Methods("PUT")
	app.Router.HandleFunc(routes.PutGuestsUri, app.deleteGuestsHandler()).Methods("DELETE")

}

func (app *App) addInvitationRoutes(){
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.GetInvitationUri, app.handleInvitationGet()).Methods("GET")
}

// Models
type PresentGuestList struct {
	Guests []PresentGuest `json:"guests"`
}

type PresentGuest struct {
	Name               string `json:"name" db:"guest_name"`
	TimeArrived        string `json:"time_arrived" db:"time_arrived"`
	AccompanyingGuests int    `json:"accompanying_guests" db:"accompanying_guests"`
}


// Guests
func (app *App) getGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentGuestList := app.getPresentGuests()

		response, _ := json.Marshal(presentGuestList)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}
}

func (app *App) updateGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (app *App) deleteGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := app.dbClient.Exec(queries.DeleteGuest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Invitation
func (app *App)handleInvitationGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
