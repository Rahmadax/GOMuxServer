package empty_seats

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/golang/mock/gomock"
	"testing"
)

func setupHandlerTests(t *testing.T) (*emptySeatsHandler, *MockEmptySeatsService) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	esServiceMock := NewMockEmptySeatsService(controller)

	esHandler := NewEmptySeatsHandler(esServiceMock, conf.Configuration{})

	return esHandler, esServiceMock
}

//func Test_countPresentGuestsHandler_Success(t *testing.T) {
//	esHandler, esServiceMock := setupHandlerTests(t)
//
//	esServiceMock.EXPECT().countPresentGuests().Return(10, nil).Times(1)
//
//	res, returnedError := esHandler.countEmptySeatsHandler()
//	assert.NoError(t, returnedError)
//	assert.Equal(t, res, 10)
//}





