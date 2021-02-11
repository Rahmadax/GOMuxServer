package guests_repo

import (
	"database/sql"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"time"
)

type guestsRepository struct {
	dbClient *sql.DB
}

func (guestsRepo *guestsRepository) UpdateArrivedGuest(name string, accompanyingGuests int) error {
	_, err := guestsRepo.dbClient.Exec(UpdateGuestArrives, accompanyingGuests, time.Now(), name)
	if err != nil {
		return err
	}

	return nil
}

func (guestsRepo *guestsRepository) GetPresentGuests() (models.PresentGuestList, error) {
	results, err := guestsRepo.dbClient.Query(GetPresentGuests)
	if err != nil {
		return models.PresentGuestList{}, err
	}

	presentGuestList := models.PresentGuestList{Guests: make([]models.PresentGuest, 0)}

	for results.Next() {
		var name string
		var accompanyingGuests int
		var timeArrived string

		err = results.Scan(&name, &accompanyingGuests, &timeArrived,)
		if err != nil {
			return models.PresentGuestList{}, err
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
		return err
	}

	return nil
}

func (guestsRepo *guestsRepository) DeleteGuest(guestName string) error {
	_, err := guestsRepo.dbClient.Exec(DeleteFromGuestList, time.Now(), guestName)
	if err != nil {
		return err
	}
	return nil
}

func (guestsRepo *guestsRepository) GetGuestList() (models.GuestList, error){
	results, err := guestsRepo.dbClient.Query(GetGuestList)
	if err != nil {
		return models.GuestList{}, err
	}

	guestList := models.GuestList{}

	for results.Next() {
		var name string
		var table, accompanyingGuests int

		err = results.Scan(&name, &table, &accompanyingGuests)
		if err != nil {
			return models.GuestList{}, err
		}

		guestList.Guests = append(guestList.Guests, models.Guest{Name: name, Table: table, AccompanyingGuests: accompanyingGuests})
	}

	return guestList, nil
}

func (guestsRepo *guestsRepository) AddToGuestList(newGuest models.Guest) error {
	_, err := guestsRepo.dbClient.Exec(InsertGuest, newGuest.Name, newGuest.Table, newGuest.AccompanyingGuests)
	if err != nil {
		return err
	}

	return nil
}

func (guestsRepo *guestsRepository) GetGuestsAtTable(tableNumber int) (int, error) {
	currentGuestsAtTable := 0
	err := guestsRepo.dbClient.QueryRow(CountExpectedGuestsAtTable, tableNumber).Scan(&currentGuestsAtTable)
	if err != nil {
		return 0, err
	}

	return currentGuestsAtTable, nil
}

func (guestsRepo *guestsRepository) DeleteFromGuestList(guestName string) error {
	_, err := guestsRepo.dbClient.Exec(DeleteFromGuestList, time.Now(), guestName)
	if err != nil {
		return err
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
		return models.FullGuestDetails{}, err
	}

	return models.FullGuestDetails{Name: name, Table: tableNumber, AccompanyingGuests: expectedGuests}, nil
}

func (guestsRepo *guestsRepository) CountPresentGuests() (int, error) {
	var currentPresentGuests int

	err := guestsRepo.dbClient.QueryRow(CountPresentGuests).Scan(&currentPresentGuests)
	if err != nil {
		return 0, err
	}

	return currentPresentGuests, nil
}

func NewGuestsRepository(dbClient *sql.DB) *guestsRepository {
	return &guestsRepository{
		dbClient: dbClient,
	}
}
