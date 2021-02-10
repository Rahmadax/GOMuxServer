package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/queries"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func (app *App) addGuestsRoutes() {
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.GetGuestsUri, app.getGuestsHandler()).Methods("GET")
	app.Router.HandleFunc(routes.PutGuestsUri, app.updateGuestsHandler()).Methods("PUT")
	app.Router.HandleFunc(routes.DeleteGuestsUri, app.deleteGuestsHandler()).Methods("DELETE")

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

type UpdateGuestRequest struct {
	AccompanyingGuests int    `json:"accompanying_guests" db:"accompanying_guests"`
}

// Guests
func (app *App) getGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentGuestList, err := app.getPresentGuests()
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(presentGuestList)
		_, _ = w.Write(response)
	}
}

func (app *App) updateGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		updateGuestReq := UpdateGuestRequest{}
		err := json.Unmarshal(body, &updateGuestReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		guestName := mux.Vars(r)["name"]
		if !isValidGuestName(guestName) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !isValidGuestNumber(updateGuestReq.AccompanyingGuests) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		storedGuestDetails, err := app.getGuest(guestName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Fewer guests than planned is always okay
		accompanyingGuestDifference := updateGuestReq.AccompanyingGuests - storedGuestDetails.AccompanyingGuests
		if accompanyingGuestDifference > 0 {
			expectedSpace, err := app.getExpectedSpaceAtTable(storedGuestDetails.Table)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Is there going to be enough space for everyone who is expected + additional newcomers?
			newExpectedSpace := expectedSpace - accompanyingGuestDifference
			if newExpectedSpace < 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		err = app.updateArrivedGuest(guestName, updateGuestReq.AccompanyingGuests)
		if err != nil {
			fmt.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(NameResponse{guestName})
		_, _ = w.Write(response)
	}
}

func (app *App) deleteGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		if !isValidGuestName(guestName) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := app.dbClient.Exec(queries.DeleteGuest, guestName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(NameResponse{Name: guestName})
		_, _ = w.Write(response)
	}
}

// Invitation
func (app *App) handleInvitationGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
