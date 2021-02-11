package guest_list

import (
	"errors"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGuestListSuccess(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	guestRepoMock := NewMockGuestsRepository(controller)
	glService := NewGuestListService(conf.Configuration{}, guestRepoMock)

	guest1 := models.Guest{Name: "Ollie", Table: 1, AccompanyingGuests: 2}
	guest2 := models.Guest{Name: "Bill", Table: 3, AccompanyingGuests: 6}
	guestList := models.GuestList{Guests: []models.Guest{guest1, guest2}}

	guestRepoMock.EXPECT().GetGuestList().Return(guestList, nil).Times(1)

	res, returnedError := glService.getGuestList()
	assert.NoError(t, returnedError)
	assert.Equal(t, res, guestList)
}

func TestGetGuestListFail_DBErrors(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	guestRepoMock := NewMockGuestsRepository(controller)
	glService := NewGuestListService(conf.Configuration{}, guestRepoMock)

	guestRepoMock.EXPECT().GetGuestList().Return(models.GuestList{}, errors.New("internal Server Error")).Times(1)

	res, returnedError := glService.getGuestList()
	assert.Error(t, returnedError)
	assert.Equal(t, res, models.GuestList{})
}
