package guests

import (
	"errors"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
)

type GuestsRepository interface {
	UpdateArrivedGuest(name string, accompanyingGuests int) error
	GetPresentGuests() (models.PresentGuestList, error)
	GetFullGuestDetails(name string) (models.FullGuestDetails, error)
	UpdateGuestLeaves(guestName string) error
	DeleteGuest(guestName string) error
	GetGuestsAtTable(tableNumber int) (int, error)
}

type guestsService struct {
	config conf.Configuration
	guestsRepo GuestsRepository
}

func (guestsService *guestsService) getPresentGuests() (models.PresentGuestList, error) {
	return guestsService.guestsRepo.GetPresentGuests()
}

func (guestsService *guestsService) guestArrives(updateGuestsReq models.UpdateGuestRequest, guestName string) error {
	storedGuestDetails, err := guestsService.guestsRepo.GetFullGuestDetails(guestName)
	if err != nil {
		return err
	}

	if storedGuestDetails.TimeArrived != nil {
		return errors.New( "guest has already arrived")
	}

	accompanyingGuestDifference := updateGuestsReq.AccompanyingGuests - storedGuestDetails.AccompanyingGuests

	// Fewer guests than planned is always okay
	if accompanyingGuestDifference > 0 {
		spaceAtTable, err := guestsService.guestsRepo.GetGuestsAtTable(storedGuestDetails.Table)
		if err != nil {
			return err
		}

		// Is there going to be enough space for everyone who is expected + additional newcomers?
		newExpectedSpace := spaceAtTable - accompanyingGuestDifference
		if newExpectedSpace < 0 {
			return errors.New(fmt.Sprintf("Not enough space expected at table. %d spaces left", spaceAtTable + storedGuestDetails.AccompanyingGuests + 1))
		}
	}

	err = guestsService.guestsRepo.UpdateArrivedGuest(guestName, updateGuestsReq.AccompanyingGuests)
	if err != nil {
		return err
	}

	return nil
}

func (guestsService *guestsService) guestLeaves(guestName string) error {
	return guestsService.guestsRepo.UpdateGuestLeaves(guestName)
}

func NewGuestsService(config conf.Configuration, guestsRepo GuestsRepository) *guestsService {
	return &guestsService{
		config: config,
		guestsRepo: guestsRepo,
	}
}