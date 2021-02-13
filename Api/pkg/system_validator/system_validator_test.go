package system_validator

import (
	"errors"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupValidatorTests(t *testing.T) *systemValidator {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tableConfig := conf.TableConfig{
		TableCapacityArray: []int{10},
		TableCapacityMap:   map[int]int{0: 10},
		TableCount:         1, TotalCapacity: 10,
	}
	systemValidator := NewSystemValidator(conf.Configuration{Tables: tableConfig})

	return systemValidator
}

func Test_ValidateGuestName_Table(t *testing.T) {
	type testStruct struct {
		guestName      string
		expectedResult error
	}

	systemValidator := setupValidatorTests(t)
	tests := []testStruct{
		{guestName: "Ted", expectedResult: nil},
		{guestName: "Jim-Jam", expectedResult: nil},
		{guestName: "D'Andre", expectedResult: nil},
		{guestName: "Dr.Dre", expectedResult: nil},
		{guestName: "Dr.Dre", expectedResult: nil},

		{guestName: "", expectedResult: errors.New("invalid guest name")},
		{guestName: "Jim Bean", expectedResult: errors.New("invalid guest name")},
		{guestName: "Jim5", expectedResult: errors.New("invalid guest name")},
		{guestName: "Jim/", expectedResult: errors.New("invalid guest name")},
		{guestName: "Jim\\", expectedResult: errors.New("invalid guest name")},

	}

	for _, test := range tests {
		returnedError := systemValidator.ValidateGuestName(test.guestName)
		if test.expectedResult != nil {
			assert.EqualError(t, returnedError, test.expectedResult.Error())
		} else {
			assert.NoError(t, returnedError)
		}
	}
}

func Test_ValidateAccompanyingGuests_Table(t *testing.T) {
	type testStruct struct {
		guestCount     int
		expectedResult error
	}

	systemValidator := setupValidatorTests(t)
	tests := []testStruct{
		{guestCount: 0, expectedResult: nil},
		{guestCount: 9, expectedResult: nil},
		{guestCount: 10, expectedResult: nil},
		{guestCount: -1, expectedResult: errors.New("guest can't have negative accompanying guests")},
	}

	for _, test := range tests {
		returnedError := systemValidator.ValidateAccompanyingGuests(test.guestCount)
		if test.expectedResult != nil {
			assert.EqualError(t, returnedError, test.expectedResult.Error())
		} else {
			assert.NoError(t, returnedError)
		}
	}
}

func Test_ValidateTableNumber_Table(t *testing.T) {
	type testStruct struct {
		tableNumber    int
		expectedResult error
	}

	systemValidator := setupValidatorTests(t)
	tests := []testStruct{
		{tableNumber: 1, expectedResult: nil},
		{tableNumber: 10, expectedResult: errors.New("there are only 1 tables")},
		{tableNumber: 0, expectedResult: errors.New("table number must be larger than 0")},
		{tableNumber: -1, expectedResult: errors.New("table number must be larger than 0")},
	}

	for _, test := range tests {
		returnedError := systemValidator.ValidateTableNumber(test.tableNumber)
		if test.expectedResult != nil {
			assert.EqualError(t, returnedError, test.expectedResult.Error())
		} else {
			assert.NoError(t, returnedError)
		}
	}
}
