package guest_list

import (
	"errors"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupServiceTests(t *testing.T) (*guestListService, *MockGuestsRepository) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	guestRepoMock := NewMockGuestsRepository(controller)
	tableConfig := conf.TableConfig{
		TableCapacityArray: []int{1, 1, 2, 3, 5, 8},
		TableCapacityMap:   map[int]int{0: 1, 1: 1, 2: 2, 3: 3, 4: 5, 5: 8},
		TableCount:         6, TotalCapacity: 20,
	}
	glService := NewGuestListService(conf.Configuration{Tables: tableConfig}, guestRepoMock)

	return glService, guestRepoMock
}

// Get guest list
func Test_GetGuestList_Success(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	guest1 := models.Guest{Name: "Ollie", Table: 1, AccompanyingGuests: 2}
	guest2 := models.Guest{Name: "Bill", Table: 3, AccompanyingGuests: 6}
	guestList := models.GuestList{Guests: []models.Guest{guest1, guest2}}

	mockGuestRepo.EXPECT().GetGuestList().Return(guestList, nil).Times(1)

	res, returnedError := glService.getGuestList()
	assert.NoError(t, returnedError)
	assert.Equal(t, res, guestList)
}

func Test_GetGuestListFail_DBErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)
	mockGuestRepo.EXPECT().GetGuestList().Return(models.GuestList{}, errors.New("internal Server Error")).Times(1)

	res, returnedError := glService.getGuestList()
	assert.Error(t, returnedError)
	assert.Equal(t, res, models.GuestList{})
}

// Add to guest list
func Test_AddToGuestListFail_Success(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	newGuest := models.Guest{
		Name:               "Alyx",
		Table:              5,
		AccompanyingGuests: 2,
	}

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(newGuest.Table).Return(2, nil).Times(1)

	mockGuestRepo.EXPECT().AddToGuestList(newGuest).Return(nil)

	res := glService.addToGuestList(newGuest)
	assert.NoError(t, res)
}

func Test_AddToGuestList_NotEnoughExpectedSpace(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	newGuest := models.Guest{
		Name:               "Alyx",
		Table:              5,
		AccompanyingGuests: 5,
	}

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(newGuest.Table).Return(2, nil).Times(1)

	res := glService.addToGuestList(newGuest)
	assert.EqualError(t, res, "Not enough space expected at table. 3 spaces left")
}

func Test_AddToGuestList_GetExpectedGuestsAtTableErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	newGuest := models.Guest{
		Name:               "Alyx",
		Table:              5,
		AccompanyingGuests: 5,
	}

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(newGuest.Table).Return(0, errors.New("something Went wrong")).Times(1)

	res := glService.addToGuestList(newGuest)
	assert.EqualError(t, res, "something Went wrong")
}

func Test_AddToGuestList_AddToGuestListErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	newGuest := models.Guest{
		Name:               "Alyx",
		Table:              5,
		AccompanyingGuests: 2,
	}

	mockGuestRepo.EXPECT().GetExpectedGuestsAtTable(newGuest.Table).Return(2, nil).Times(1)

	mockGuestRepo.EXPECT().AddToGuestList(newGuest).Return(errors.New("something Went wrong"))
	res := glService.addToGuestList(newGuest)
	assert.EqualError(t, res, "something Went wrong")
}

// Remove from guest list
func Test_RemoveFromGuestList_Success(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	mockGuestRepo.EXPECT().GetFullGuestDetails("GordonFreeman").Return(
		models.FullGuestDetails{Name: "GordonFreeman", Table: 10, AccompanyingGuests: 5}, nil,
	).Times(1)

	mockGuestRepo.EXPECT().DeleteFromGuestList("GordonFreeman").Return(
		nil,
	).Times(1)

	res := glService.removeFromGuestList("GordonFreeman")
	assert.NoError(t, res)
}

func Test_RemoveFromGuestList_GuestIsAlreadyCheckedIn(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	arrivalTime := "2019-10-19T17:12:30.174"
	mockGuestRepo.EXPECT().GetFullGuestDetails("GordonFreeman").Return(
		models.FullGuestDetails{Name: "GordonFreeman", Table: 10, AccompanyingGuests: 5, TimeArrived: &arrivalTime}, nil,
	).Times(1)

	res := glService.removeFromGuestList("GordonFreeman")
	assert.EqualError(t, res, "A guest that has already arrived cannot be removed from the guest list")
}

func Test_RemoveFromGuestList_GetFullDetailsDbErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	mockGuestRepo.EXPECT().GetFullGuestDetails("GordonFreeman").Return(
		models.FullGuestDetails{}, errors.New("something went wrong"),
	).Times(1)

	res := glService.removeFromGuestList("GordonFreeman")
	assert.EqualError(t, res, "something went wrong")
}

func Test_RemoveFromGuestList_DeleteFromGuestListDbErrors(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	mockGuestRepo.EXPECT().GetFullGuestDetails("GordonFreeman").Return(
		models.FullGuestDetails{Name: "GordonFreeman", Table: 10, AccompanyingGuests: 5}, nil,
	).Times(1)

	mockGuestRepo.EXPECT().DeleteFromGuestList("GordonFreeman").Return(
		errors.New("something went wrong"),
	).Times(1)

	res := glService.removeFromGuestList("GordonFreeman")
	assert.EqualError(t, res, "something went wrong")
}
