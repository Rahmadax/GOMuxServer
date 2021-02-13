package guest_list

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/system_validator"
	"github.com/golang/mock/gomock"
	"testing"
)

func setupHandlerTests(t *testing.T) (*guestListHandler, *MockGuestListService) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockGlService := NewMockGuestListService(controller)
	tableConfig := conf.TableConfig{
		TableCapacityArray: []int{1, 1, 2, 3, 5, 8},
		TableCapacityMap:   map[int]int{0: 1, 1: 1, 2: 2, 3: 3, 4: 5, 5: 8},
		TableCount:         6, TotalCapacity: 20,
	}
	systemValidator := system_validator.NewSystemValidator(conf.Configuration{Tables: tableConfig})

	glHandler := NewGuestListHandler(mockGlService, systemValidator)

	return glHandler, mockGlService
}

//// Get guest list
//func TestGetGuestListHandlerSuccess(t *testing.T) {
//	glHandler, mockGlService := setupHandlerTests(t)
//
//	guest1 := models.Guest{Name: "Ollie", Table: 1, AccompanyingGuests: 2}
//	guest2 := models.Guest{Name: "Bill", Table: 3, AccompanyingGuests: 6}
//	guestList := models.GuestList{Guests: []models.Guest{guest1, guest2}}
//
//	mockGlService.EXPECT().getGuestList().Return(guestList, nil).Times(1)
//
//	res := glHandler.getGuestListHandler()
//
//	assert.Equal(t, res, guestList)
//}
