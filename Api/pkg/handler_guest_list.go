package pkg

import (
	"encoding/json"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/queries"
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
	Name               string `json:"name"`
	Table              int    `json:"table"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

type NameResponse struct {
	Name string `json:"name"`
}

func (app *App) getGuestListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := app.dbClient.Query(queries.GetGuestList)
		if err != nil {
			panic("AHHHH")
		}

		guestList := GuestList{}

		for results.Next() {
			thisGuest := Guest{}

			err = results.Scan(thisGuest)
			if err != nil {
				panic(err.Error())
			}

			guestList.Guests = append(guestList.Guests, thisGuest)
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

		currentGuestsAtTable := 0
		_ = app.dbClient.QueryRow(queries.CountGuestsAtTable).Scan(currentGuestsAtTable)

		remainingSpace := app.Config.Tables.TableCapacity - currentGuestsAtTable
		if newGuest.AccompanyingGuests+1 > remainingSpace {
			w.WriteHeader(http.StatusBadRequest)
		}

		_, _ = app.dbClient.Exec(queries.InsertGuest, newGuest.Name, newGuest.Table, newGuest.AccompanyingGuests)

		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(newGuest.Name)
		_, _ = w.Write(response)
	}
}

func (app *App) guestListDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		newGuest := Guest{}
		_ = json.Unmarshal(body, &newGuest)

		_, _ = app.dbClient.Exec(queries.DeleteGuest, newGuest.Name)

		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(NameResponse{Name: newGuest.Name})
		w.Write(response)
	}
}
