package guests

import (
	"errors"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupServiceTests(t *testing.T) (*guestsService, *MockGuestsRepository) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	guestRepoMock := NewMockGuestsRepository(controller)
	tableConfig := conf.TableConfig{
		TableCapacityArray: []int{1, 1, 2, 3, 5, 8},
		TableCapacityMap:   map[int]int{0: 1, 1: 1, 2: 2, 3: 3, 4: 5, 5: 8},
		TableCount:         6, TotalCapacity: 20,
	}
	glService := NewGuestsService(conf.Configuration{Tables: tableConfig}, guestRepoMock)

	return glService, guestRepoMock
}

// Get present guest list
func TestGetPresentGuestsSuccess(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guest1 := models.PresentGuest{Name: "Ollie", AccompanyingGuests: 2}
	guest2 := models.PresentGuest{Name: "Bill",  AccompanyingGuests: 6}
	presentGuestList := models.PresentGuestList{Guests: []models.PresentGuest{guest1, guest2}}

	mockGuestRepo.EXPECT().GetPresentGuests().Return(presentGuestList, nil).Times(1)

	res, returnedError := glService.getPresentGuests()
	assert.NoError(t, returnedError)
	assert.Equal(t, res, presentGuestList)
}

func TestGetPresentGuestsErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 5, AccompanyingGuests: 4}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(fullDetails.Table).Return(2, nil).Times(1)

	mockGuestRepo.EXPECT().UpdateGuestLeaves(fullDetails.Name).Return( nil).Times(1)

	returnedError := glService.guestLeaves(guestName)
	assert.NoError(t, returnedError)
}

// Guest Arrives
func TestGuestArrivesGuestsSuccess(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"
	newAccompanyingGuestCount := 4

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 6, AccompanyingGuests: 4}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(fullDetails.Table).Return(2, nil).Times(1)

	mockGuestRepo.EXPECT().UpdateArrivedGuest(fullDetails.Name, newAccompanyingGuestCount).Return( nil).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.NoError(t, returnedError)
}

func TestGuestArrivesGuests__FewerAccompanyingSuccess(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"
	newAccompanyingGuestCount := 1

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 6, AccompanyingGuests: 4}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	mockGuestRepo.EXPECT().UpdateArrivedGuest(fullDetails.Name, newAccompanyingGuestCount).Return( nil).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.NoError(t, returnedError)
}

func TestGuestArrivesGuests__MoreAccompanyingSuccess(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"
	newAccompanyingGuestCount := 5

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 6, AccompanyingGuests: 4}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(fullDetails.Table).Return(2, nil).Times(1)

	mockGuestRepo.EXPECT().UpdateArrivedGuest(fullDetails.Name, newAccompanyingGuestCount).Return( nil).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.NoError(t, returnedError)
}

func TestGuestArrivesGuests__MoreAccompanyingFailure(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)
	// Total table space = 8

	guestName := "Tan"
	newAccompanyingGuestCount := 50

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 6, AccompanyingGuests: 3}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(fullDetails.Table).Return(2, nil).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.EqualError(t, returnedError, "Not enough space expected at table. 6 spaces left")
}

func TestGuestArrivesGuests__GuestIsAlreadyHere(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"
	newAccompanyingGuestCount := 5
	timeArrived := "2012.12.42:00:30:39"

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 6, AccompanyingGuests: 3, TimeArrived: &timeArrived}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.EqualError(t, returnedError, "guest has already arrived")
}

func TestGuestArrivesGuests__GetFullGuestDetailsErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"
	newAccompanyingGuestCount := 5

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(models.FullGuestDetails{}, errors.New("something went wrong")).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.EqualError(t, returnedError, "something went wrong")
}

func TestGuestArrivesGuests__GetExpectedGuestsAtTableErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"
	newAccompanyingGuestCount := 5

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 6, AccompanyingGuests: 3}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(fullDetails.Table).Return(2, nil).Times(1)

	mockGuestRepo.EXPECT().UpdateArrivedGuest(fullDetails.Name, newAccompanyingGuestCount).Return(errors.New("something went wrong")).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.EqualError(t, returnedError, "something went wrong")
}

func TestGuestArrivesGuests__UpdateArrivedGuestErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guestName := "Tan"
	newAccompanyingGuestCount := 5

	fullDetails := models.FullGuestDetails{Name: guestName, Table: 6, AccompanyingGuests: 3}

	mockGuestRepo.EXPECT().GetFullGuestDetails(guestName).Return(fullDetails, nil).Times(1)

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(fullDetails.Table).Return(0, errors.New("something went wrong")).Times(1)

	returnedError := glService.guestArrives(newAccompanyingGuestCount, guestName)
	assert.EqualError(t, returnedError, "something went wrong")
}

// Guest Leaves
func TestGuestLeavesSuccess(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)
	guestName := "Fred"

	mockGuestRepo.EXPECT().UpdateGuestLeaves(guestName).Return(nil).Times(1)

	returnedError := glService.guestLeaves(guestName)
	assert.NoError(t, returnedError)
}

func TestGuestLeavesErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)
	guestName := "Fred"

	mockGuestRepo.EXPECT().UpdateGuestLeaves(guestName).Return(errors.New("something went wrong")).Times(1)

	returnedError := glService.guestLeaves(guestName)
	assert.EqualError(t, returnedError, "something went wrong")
}


