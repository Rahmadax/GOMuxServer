package pkg

import (
	"encoding/json"
	"fmt"
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
		results, err := app.dbClient.Query(queries.GetGuestList)
		if err != nil {
			panic(err)
		}

		guestList := GuestList{}

		for results.Next() {
			var name string
			var table int
			var accompanyingGuests int

			err = results.Scan(&name, &table, &accompanyingGuests)
			if err != nil {
				panic(err.Error())
			}

			guestList.Guests = append(guestList.Guests, Guest{name, table, accompanyingGuests})
		}

		response, _ := json.Marshal(guestList)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(response)
	}
}

func (app *App) postGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, _ := ioutil.ReadAll(r.Body)
		newGuest := Guest{}
		_ = json.Unmarshal(body, &newGuest)
		newGuest.Name = mux.Vars(r)["name"]

		currentGuestsAtTable := 0
		err := app.dbClient.QueryRow(queries.CountGuestsAtTable, newGuest.Table).Scan(&currentGuestsAtTable)
		if err != nil {
			panic(err)
		}

		remainingSpace := app.Config.Tables.TableCapacity - currentGuestsAtTable
		if newGuest.AccompanyingGuests+1 > remainingSpace {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(fmt.Sprintf("Remaining Sapce: %s, new Guests: %s", remainingSpace, currentGuestsAtTable))
			return
		}

		_, err = app.dbClient.Exec(queries.InsertGuest, newGuest.Name, newGuest.Table, newGuest.AccompanyingGuests)
		if err != nil {
			panic(err)
		}

		response, _ := json.Marshal(NameResponse{newGuest.Name})
		_, _ = w.Write(response)
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
