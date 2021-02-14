package guests

import (
	"errors"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/system_validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupHandlerTests(t *testing.T) (*guestsHandler, *MockGuestsService) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockGlService := NewMockGuestsService(controller)
	tableConfig := conf.TableConfig{
		TableCapacityArray: []int{1, 1, 2, 3, 5, 8},
		TableCapacityMap:   map[int]int{0: 1, 1: 1, 2: 2, 3: 3, 4: 5, 5: 8},
		TableCount:         6, TotalCapacity: 20,
	}
	systemValidator := system_validator.NewSystemValidator(conf.Configuration{Tables: tableConfig})

	glHandler := newGuestListHandler(mockGlService, systemValidator)

	return glHandler, mockGlService
}

func Test_GuestArrives__Success(t *testing.T) {
	glHandler, mockGlService := setupHandlerTests(t)

	updateRequest := models.UpdateGuestRequest{AccompanyingGuests: 5}
	guestName := "Ted"
	expectedRes := models.NameResponse{Name: "Ted"}

	mockGlService.EXPECT().guestArrives(updateRequest.AccompanyingGuests, guestName).Return(nil).Times(1)

	res, err := glHandler.guestArrives(updateRequest, guestName)
	assert.NoError(t, err)
	assert.Equal(t, res, expectedRes)
}

func Test_GuestArrives_ValidationFailureNegative(t *testing.T) {
	glHandler, _ := setupHandlerTests(t)

	updateRequest := models.UpdateGuestRequest{AccompanyingGuests: -1}
	guestName := "Ted"
	expectedRes := models.NameResponse{}

	res, err := glHandler.guestArrives(updateRequest, guestName)
	assert.EqualError(t, err, "guest can't have negative accompanying guests")
	assert.Equal(t, res, expectedRes)
}

func Test_GuestArrives_ServiceError(t *testing.T) {
	glHandler, mockGlService := setupHandlerTests(t)

	updateRequest := models.UpdateGuestRequest{AccompanyingGuests: 5}
	guestName := "Ted"

	mockGlService.EXPECT().guestArrives(updateRequest.AccompanyingGuests, guestName).Return(errors.New("something went wrong")).Times(1)

	res, err := glHandler.guestArrives(updateRequest, guestName)
	assert.EqualError(t, err, "something went wrong")
	assert.Equal(t, res, models.NameResponse{})
}