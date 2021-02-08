package pkg

import (
	"encoding/json"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/queries"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
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
	Name               string `json:"name" db:"guest_name"`
	Table              int    `json:"table" db:"table_number"`
	AccompanyingGuests int    `json:"accompanying_guests" db:"accompanying_guests"`
}

type NameResponse struct {
	Name string `json:"name"`
}

func (app *App) getGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestList, err := app.getGuestList()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(guestList)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func (app *App) postGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, _ := ioutil.ReadAll(r.Body)
		newGuest := Guest{}
		_ = json.Unmarshal(body, &newGuest)
		newGuest.Name = mux.Vars(r)["name"]

		if ok := newGuest.validate(app.Config.Tables.TableCount); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if ok, err := app.isSpaceAtTable(newGuest); err != nil || !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := app.insertGuest(newGuest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(NameResponse{newGuest.Name})

		w.Write(response)
		w.WriteHeader(http.StatusOK)
	}
}

func (app *App) guestListDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]

		_, err := app.dbClient.Exec(queries.DeleteGuest, guestName)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(NameResponse{Name: guestName})
		w.Write(response)
	}
}
