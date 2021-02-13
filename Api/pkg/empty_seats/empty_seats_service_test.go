package empty_seats

import (
	"errors"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupServiceTests(t *testing.T) (*emptySeatsService, *MockGuestsRepository) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	guestRepoMock := NewMockGuestsRepository(controller)

	esService := NewEmptySeatsService(conf.Configuration{}, guestRepoMock)

	return esService, guestRepoMock
}

func Test_countPresentGuests_Success(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	mockGuestRepo.EXPECT().CountPresentGuests().Return(10, nil).Times(1)

	res, returnedError := glService.countPresentGuests()
	assert.NoError(t, returnedError)
	assert.Equal(t, res, 10)
}

func Test_countPresentGuests_Failure(t *testing.T) {
	glService, mockGuestRepo := setupServiceTests(t)

	mockGuestRepo.EXPECT().CountPresentGuests().Return(0, errors.New("something went wrong")).Times(1)

	res, returnedError := glService.countPresentGuests()
	assert.EqualError(t, returnedError,"something went wrong")
	assert.Equal(t, res, 0)
}





