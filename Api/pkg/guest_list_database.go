package pkg

import (
	"github.com/Rahmadax/GOMuxServer/Api/pkg/queries"
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

func (app *App) isSpaceAtTable(newGuest Guest) (bool, error) {
	currentGuestsAtTable := 0
	err := app.dbClient.QueryRow(queries.CountGuestsAtTable, newGuest.Table).Scan(&currentGuestsAtTable)
	if err != nil {
		return false, err
	}

	remainingSpace := app.Config.Tables.TableCapacity - currentGuestsAtTable
	if newGuest.AccompanyingGuests+1 > remainingSpace {
		return false, nil
	}

	return true, nil
}

func (app *App) insertGuest(newGuest Guest) error {
	_, err := app.dbClient.Exec(queries.InsertGuest, newGuest.Name, newGuest.Table, newGuest.AccompanyingGuests)
	if err != nil {
		return err
	}
	return nil
}