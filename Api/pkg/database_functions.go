package pkg

import (
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/queries"
	"time"
)

func (app *App) getGuestList() (GuestList, error){
	results, err := app.dbClient.Query(queries.GetGuestList)
	if err != nil {
		return GuestList{}, err
	}

	guestList := GuestList{}

	for results.Next() {
		var name string
		var table, accompanyingGuests int

		err = results.Scan(&name, &table, &accompanyingGuests)
		if err != nil {
			return GuestList{}, err
		}

		guestList.Guests = append(guestList.Guests, Guest{name, table, accompanyingGuests})
	}
	return guestList, nil
}

func (app *App) deleteAll() {
	_, err := app.dbClient.Exec(queries.DeleteAll)
	fmt.Println(err)
}

func (app *App) getExpectedSpaceAtTable(tableNumber int) (int, error) {
	currentGuestsAtTable := 0
	err := app.dbClient.QueryRow(queries.CountExpectedGuestsAtTable, tableNumber).Scan(&currentGuestsAtTable)
	if err != nil {
		return 0, err
	}

	return app.Config.Tables.TableCapacity - currentGuestsAtTable, nil
}

func (app *App) updateArrivedGuest(name string, accompanyingGuests int) error {
	_, err := app.dbClient.Exec(queries.UpdateArrivedGuest, accompanyingGuests, time.Now(), name)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) insertGuest(newGuest Guest) error {
	_, err := app.dbClient.Exec(queries.InsertGuest, newGuest.Name, newGuest.Table, newGuest.AccompanyingGuests)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) getPresentGuests() (PresentGuestList, error) {
	results, err := app.dbClient.Query(queries.GetPresentGuests)
	if err != nil {
		return PresentGuestList{}, err
	}

	presentGuestList := PresentGuestList{Guests:make([]PresentGuest, 0)}

	for results.Next() {
		var name string
		var accompanyingGuests int
		var timeArrived string

		err = results.Scan(&name, &accompanyingGuests, &timeArrived,)
		if err != nil {
			return PresentGuestList{}, err
		}

		presentGuestList.Guests = append(presentGuestList.Guests,
			PresentGuest{
				Name: name,
				AccompanyingGuests: accompanyingGuests,
				TimeArrived: timeArrived,
			})
	}

	return presentGuestList, nil
}

func (app *App) getFullGuestDetails(name string) (FullGuestDetails, error){
	var tableNumber int
	var expectedGuests int
	var timeArrived *string
	var timeLeft *string

	err := app.dbClient.QueryRow(queries.GetGuest, name).Scan(&tableNumber, &expectedGuests, &timeArrived, &timeLeft)
	if err != nil {
		return FullGuestDetails{}, err
	}

	return FullGuestDetails{Name: name, Table: tableNumber, AccompanyingGuests: expectedGuests}, nil
}

func (app *App) countPresentGuests() (int, error) {
	var currentPresentGuests int

	err := app.dbClient.QueryRow(queries.CountPresentGuests).Scan(&currentPresentGuests)
	if err != nil {
		return 0, err
	}

	return currentPresentGuests, nil
}