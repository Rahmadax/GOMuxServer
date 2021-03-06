package guests_repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"time"
)

type guestsRepository struct {
	dbClient *sql.DB
}

func (guestsRepo *guestsRepository) UpdateArrivedGuest(name string, accompanyingGuests int) error {
	_, err := guestsRepo.dbClient.Exec(UpdateGuestArrives, accompanyingGuests, time.Now(), name)
	if err != nil {
		fmt.Println(err) // This should feed to a local kibana or similar, but ran out of time
		return errors.New("something went wrong")
	}

	return nil
}

func (guestsRepo *guestsRepository) GetPresentGuests() (models.PresentGuestList, error) {
	results, err := guestsRepo.dbClient.Query(GetPresentGuests)
	if err != nil {
		fmt.Println(err)
		return models.PresentGuestList{}, errors.New("something went wrong")
	}

	presentGuestList := models.PresentGuestList{Guests: make([]models.PresentGuest, 0)}

	for results.Next() {
		var name string
		var accompanyingGuests int
		var timeArrived string

		err = results.Scan(&name, &accompanyingGuests, &timeArrived,)
		if err != nil {
			fmt.Println(err)
			return models.PresentGuestList{}, errors.New("something went wrong")
		}

		presentGuestList.Guests = append(presentGuestList.Guests,
			models.PresentGuest{
				Name: name,
				AccompanyingGuests: accompanyingGuests,
				TimeArrived: timeArrived,
			})
	}

	return presentGuestList, nil
}

func (guestsRepo *guestsRepository) UpdateGuestLeaves(guestName string) error {
	_, err := guestsRepo.dbClient.Exec(UpdateGuestLeaves, time.Now(), guestName)
	if err != nil {
		fmt.Println(err)
		return errors.New("something went wrong")	}

	return nil
}

func (guestsRepo *guestsRepository) DeleteGuest(guestName string) error {
	_, err := guestsRepo.dbClient.Exec(DeleteFromGuestList, time.Now(), guestName)
	if err != nil {
		fmt.Println(err)
		return errors.New("something went wrong")
	}
	return nil
}

func (guestsRepo *guestsRepository) GetGuestList() (models.GuestList, error){
	results, err := guestsRepo.dbClient.Query(GetGuestList)
	if err != nil {
		fmt.Println(err)
		return models.GuestList{}, errors.New("something went wrong")
	}

	guestList := models.GuestList{}

	for results.Next() {
		var name string
		var table, accompanyingGuests int

		err = results.Scan(&name, &table, &accompanyingGuests)
		if err != nil {
			fmt.Println(err)
			return models.GuestList{}, errors.New("something went wrong")
		}

		guestList.Guests = append(guestList.Guests, models.Guest{Name: name, Table: table, AccompanyingGuests: accompanyingGuests})
	}

	return guestList, nil
}

func (guestsRepo *guestsRepository) AddToGuestList(newGuest models.Guest) error {
	_, err := guestsRepo.dbClient.Exec(InsertGuest, newGuest.Name, newGuest.Table, newGuest.AccompanyingGuests)
	if err != nil {
		fmt.Println(err)
		return errors.New("something went wrong")
	}

	return nil
}

func (guestsRepo *guestsRepository) GetExpectedGuestsAtTable(tableNumber int) (int, error) {
	currentGuestsAtTable := 0
	err := guestsRepo.dbClient.QueryRow(CountExpectedGuestsAtTable, tableNumber).Scan(&currentGuestsAtTable)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("something went wrong")
	}

	return currentGuestsAtTable, nil
}

func (guestsRepo *guestsRepository) DeleteFromGuestList(guestName string) error {
	_, err := guestsRepo.dbClient.Exec(DeleteFromGuestList, time.Now(), guestName)
	if err != nil {
		fmt.Println(err)
		return errors.New("something went wrong")
	}

	return nil
}

func (guestsRepo *guestsRepository) GetFullGuestDetails(name string) (models.FullGuestDetails, error) {
	var tableNumber int
	var expectedGuests int
	var timeArrived *string
	var timeLeft *string

	err := guestsRepo.dbClient.QueryRow(GetGuestFullDetails, name).Scan(&tableNumber, &expectedGuests, &timeArrived, &timeLeft)
	if err != nil {
		fmt.Println(err)
		return models.FullGuestDetails{}, errors.New("something went wrong")
	}

	return models.FullGuestDetails{Name: name, Table: tableNumber, AccompanyingGuests: expectedGuests}, nil
}

func (guestsRepo *guestsRepository) CountPresentGuests() (int, error) {
	var currentPresentGuests int

	err := guestsRepo.dbClient.QueryRow(CountPresentGuests).Scan(&currentPresentGuests)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("something went wrong")
	}

	return currentPresentGuests, nil
}

func NewGuestsRepository(dbClient *sql.DB) *guestsRepository {
	return &guestsRepository{
		dbClient: dbClient,
	}
}
