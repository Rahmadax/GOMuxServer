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

func (app *App) addGuestListRoutes() {
	routes := app.Config.Routes

	app.Router.HandleFunc(routes.GetGuestListUri, app.getGuestListHandler()).Methods("GET")
	app.Router.HandleFunc(routes.PostGuestListUri, app.postGuestListHandler()).Methods("POST")
	app.Router.HandleFunc(routes.DeleteGuestListUri, app.guestListDeleteHandler()).Methods("DELETE")
}

type GuestList struct {
	Guests []Guest `json:"guests"`
}

type Guest struct {
	Name               string  `json:"name" db:"guest_name"`
	Table              int     `json:"table" db:"table_number"`
	AccompanyingGuests int     `json:"accompanying_guests" db:"accompanying_guests"`
}

type FullGuestDetails struct {
	Name               string  `db:"guest_name"`
	Table              int     `db:"table_number"`
	AccompanyingGuests int     `db:"accompanying_guests"`
	TimeArrived        *string `db:"time_arrived"`
	TimeLeft           *string `db:"time_left"`
}

type NameResponse struct {
	Name string `json:"name"`
}

func (app *App) getGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestList, err := app.getGuestList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, _ := json.Marshal(guestList)
		_, _ = w.Write(response)
	}
}

func (app *App) postGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, _ := ioutil.ReadAll(r.Body)
		newGuest := Guest{}
		err := json.Unmarshal(body, &newGuest)
		if err != nil {
			handleResponseError(http.StatusBadRequest, "Invalid Request Body", w)
			return
		}

		newGuest.Name = mux.Vars(r)["name"]
		if ok := newGuest.validate(app.Config.Tables.TableCount); !ok {
			handleResponseError(http.StatusBadRequest, fmt.Sprintf("Invalid guest name: %s", newGuest.Name), w)
			return
		}

		expectedSpace, err := app.getExpectedSpaceAtTable(newGuest.Table)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if expectedSpace-newGuest.AccompanyingGuests-1 < 0 {
			handleResponseError(http.StatusBadRequest, fmt.Sprintf("Not enough space expected at table. %d spaces left", expectedSpace), w)
			return
		}

		if err := app.addGuestToGuestList(newGuest); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		response, _ := json.Marshal(NameResponse{newGuest.Name})
		_, _ = w.Write(response)
	}
}

func (app *App) guestListDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		if !isValidGuestName(guestName) {
			handleResponseError(http.StatusBadRequest, fmt.Sprintf("Invalid guest name: %s", guestName), w)
			return
		}

		guestDetails := FullGuestDetails{}
		err := app.dbClient.QueryRow(queries.GetGuestFullDetails, guestName).Scan(&guestDetails)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if guestDetails.TimeArrived != nil {
			handleResponseError(http.StatusBadRequest, fmt.Sprintf("A guest that has already arrived cannot be removed from the guest list"), w)
			return
		}

		_, err = app.dbClient.Exec(queries.DeleteFromGuestList, time.Now(), guestName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, _ := json.Marshal(NameResponse{Name: guestName})
		_, _ = w.Write(response)
	}
}
