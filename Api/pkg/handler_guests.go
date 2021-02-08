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

type PresentGuestList struct {
	Guests []PresentGuest `json:"guests"`
}

type PresentGuest struct {
	Name               string `json:"name" db:"guest_name"`
	TimeArrived        string `json:"time_arrived" db:"time_arrived"`
	AccompanyingGuests int    `json:"accompanying_guests" db:"accompanying_guests"`
}

func (app *App) getGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := app.dbClient.Query(queries.CountPresentGuests)
		if err != nil {
			panic(err)
		}

		guestList := PresentGuestList{}

		for results.Next() {
			var name string
			var accompanyingGuests int
			var timeArrived string

			err = results.Scan(&name, &accompanyingGuests, &timeArrived)
			if err != nil {
				panic(err.Error())
			}

			guestList.Guests = append(guestList.Guests, PresentGuest{name, timeArrived, accompanyingGuests})
		}

		response, _ := json.Marshal(guestList)

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

		}

		w.WriteHeader(http.StatusOK)
	}
}
