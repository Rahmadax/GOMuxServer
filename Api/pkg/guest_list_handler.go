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

	app.Router.HandleFunc("/delete_all", app.deleteAllHandler()).Methods("DELETE")

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

func (app *App) deleteAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.deleteAll()
	}
}


func (app *App) getGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestList, err := app.getGuestList()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newGuest.Name = mux.Vars(r)["name"]

		if ok := newGuest.validate(app.Config.Tables.TableCount); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		expectedSpace, err := app.getExpectedSpaceAtTable(newGuest.Table);
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if expectedSpace < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := app.insertGuest(newGuest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(NameResponse{newGuest.Name})

		_, _ = w.Write(response)
	}
}

func (app *App) guestListDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guestName := mux.Vars(r)["name"]
		if !isValidGuestName(guestName) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := app.dbClient.Exec(queries.DeleteGuest, guestName)
		if err != nil {
			panic(err)
		}

		response, _ := json.Marshal(NameResponse{Name: guestName})
		_, _ = w.Write(response)
	}
}
