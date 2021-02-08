package pkg

import (
	"encoding/json"
	"net/http"
)

func (app *App) addSeatsRoutes(){
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.CountEmptySeatsUri, app.handleCountEmptySeats()).Methods("GET")
}

type SeatsEmptyResponse struct {
	SeatsEmpty int `json:"seats_empty"`
}

func (app *App) handleCountEmptySeats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentGuestCount := app.countPresentGuests()

		maximumSeats := app.Config.Tables.TableCount * app.Config.Tables.TableCapacity

		response, _ := json.Marshal(SeatsEmptyResponse{maximumSeats - presentGuestCount})
		_, _ = w.Write(response)
	}
}