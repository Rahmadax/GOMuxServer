package guest_list

import (
	"errors"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
)

type GuestsRepository interface {
	GetFullGuestDetails(name string) (models.FullGuestDetails, error)
	GetGuestList() (models.GuestList, error)
	AddToGuestList(newGuest models.Guest) error
	DeleteFromGuestList(guestName string) error
	GetExpectedGuestsAtTable(tableNumber int) (int, error)
}

type guestListService struct {
	config conf.Configuration
	guestsRepo GuestsRepository
}

func (glService guestListService) getGuestList() (models.GuestList, error) {
	guestList, err := glService.guestsRepo.GetGuestList()
	if err != nil {
		return models.GuestList{}, err
	}

	return guestList, nil
}

func (glService guestListService) addToGuestList(newGuest models.Guest) error {
	currentGuestsAtTable, err := glService.guestsRepo.GetExpectedGuestsAtTable(newGuest.Table)
	if err != nil {
		return err
	}

	expectedSpace := glService.config.Tables.TableCapacityMap[newGuest.Table-1] - currentGuestsAtTable
	if expectedSpace-newGuest.AccompanyingGuests-1 < 0 {
		return errors.New(fmt.Sprintf("Not enough space expected at table. %d spaces left", expectedSpace))
	}

	if err := glService.guestsRepo.AddToGuestList(newGuest); err != nil {
		return err
	}

	return nil
}

func (glService guestListService) removeFromGuestList(guestName string) error {
	guestDetails, err := glService.guestsRepo.GetFullGuestDetails(guestName)
	if err != nil {
		return err
	}

	if guestDetails.TimeArrived != nil {
		return errors.New(fmt.Sprintf("A guest that has already arrived cannot be removed from the guest list"))
	}

	err = glService.guestsRepo.DeleteFromGuestList(guestName)
	if err != nil {
		return err
	}

	return nil
}

func NewGuestListService(config conf.Configuration, guestsRepo GuestsRepository) *guestListService {
	return &guestListService{
		config: config,
		guestsRepo: guestsRepo,
	}
}

