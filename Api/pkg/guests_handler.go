package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/queries"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

func (app *App) addGuestsRoutes() {
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.GetGuestsUri, app.getGuestsHandler()).Methods("GET")
	app.Router.HandleFunc(routes.PutGuestsUri, app.guestArrivesHandler()).Methods("PUT")
	app.Router.HandleFunc(routes.DeleteGuestsUri, app.guestLeavesHandler()).Methods("DELETE")

	app.Router.HandleFunc(routes.GetInvitationUri, app.handleInvitationGet()).Methods("GET")
}

// Models
type PresentGuestList struct {
	Guests []PresentGuest `json:"guests"`
}

type PresentGuest struct {
	Name               string  `json:"name" db:"guest_name"`
	AccompanyingGuests int     `json:"accompanying_guests" db:"accompanying_guests"`
	TimeArrived        string  `json:"time_arrived" db:"time_arrived"`
	TimeLeft           *string `json:",omitempty" db:"time_left"`
}

type UpdateGuestRequest struct {
	AccompanyingGuests int `json:"accompanying_guests" db:"accompanying_guests"`
}

// Guests
func (app *App) getGuestsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentGuestList, err := app.getPresentGuests()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, _ := json.Marshal(presentGuestList)
		_, _ = w.Write(response)
	}
}

func (app *App) guestArrivesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		updateGuestReq := UpdateGuestRequest{}
		err := json.Unmarshal(body, &updateGuestReq)
		if err != nil {
			handleResponseError(http.StatusBadRequest, err.Error(), w)
			return
		}

		guestName := mux.Vars(r)["name"]
		err = validateUpdateGuestRequest(guestName, updateGuestReq)
		if err != nil {
			handleResponseError(http.StatusBadRequest, err.Error(), w)
			return
		}

		storedGuestDetails, err := app.getFullGuestDetails(guestName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if storedGuestDetails.TimeArrived != nil {
			handleResponseError(http.StatusBadRequest, "Guest has already arrived", w)
			return
		}

		accompanyingGuestDifference := updateGuestReq.AccompanyingGuests - storedGuestDetails.AccompanyingGuests

		// Fewer guests than planned is always okay
		if accompanyingGuestDifference > 0 {
			expectedSpace, err := app.getExpectedSpaceAtTable(storedGuestDetails.Table)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Is there going to be enough space for everyone who is expected + additional newcomers?
			newExpectedSpace := expectedSpace - accompanyingGuestDifference
			if newExpectedSpace < 0 {
				handleResponseError(http.StatusBadRequest, fmt.Sprintf("Not enough space expected at table. %d spaces left", expectedSpace + storedGuestDetails.AccompanyingGuests + 1), w)
				return
			}
		}

		err = app.updateArrivedGuest(guestName, updateGuestReq.AccompanyingGuests)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, _ := json.Marshal(NameResponse{guestName})
		_, _ = w.Write(response)
	}
}

func (app *App) guestLeavesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		if !isValidGuestName(guestName) {
			handleResponseError(http.StatusBadRequest, fmt.Sprintf("Invalid guest name: %s", guestName), w)
			return
		}

		_, err := app.dbClient.Exec(queries.GuestLeaves, time.Now(), guestName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
