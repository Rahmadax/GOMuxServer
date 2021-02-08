package pkg

import "github.com/Rahmadax/GOMuxServer/Api/pkg/queries"

func (app *App) getPresentGuests() PresentGuestList {
	results, err := app.dbClient.Query(queries.GetPresentGuests)
	if err != nil {
		panic(err)
	}

	presentGuestList := PresentGuestList{}

	for results.Next() {
		var name string
		var accompanyingGuests int
		var timeArrived string

		err = results.Scan(&name, &accompanyingGuests, &timeArrived)
		if err != nil {
			panic(err.Error())
		}

		presentGuestList.Guests = append(presentGuestList.Guests,
			PresentGuest{
				Name: name,
				TimeArrived: timeArrived,
				AccompanyingGuests: accompanyingGuests,
			})
	}

	return presentGuestList
}