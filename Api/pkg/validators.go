package pkg

import (
	"errors"
	"regexp"
)

func (newGuest Guest) validate(tableCount int) bool {
	validNameChan := make(chan bool)
	validGuestNumChan := make(chan bool)
	validTableNumChan := make(chan bool)

	go func() {
		validNameChan <- isValidGuestName(newGuest.Name)
	}()
	go func() {
		validGuestNumChan <- isValidGuestNumber(newGuest.AccompanyingGuests)
	}()
	go func() {
		validTableNumChan <- isValidTableNumber(newGuest.Table, tableCount)
	}()

	return <-validNameChan && <-validGuestNumChan && <-validTableNumChan
}

func validateUpdateGuestRequest(guestName string, updateGuestReq UpdateGuestRequest) error {
	errString := ""
	if !isValidGuestName(guestName) {
		errString +=" invalid guest name "
	}

	if !isValidGuestNumber(updateGuestReq.AccompanyingGuests) {
		errString +="invalid number of guests"
	}

	if errString != "" {
		return errors.New(errString)
	}

	return nil
}

// general
func isValidGuestName(name string) bool {
	matched, err := regexp.MatchString(`[a-z,.'-]+`, name)
	if err != nil {
		return false
	}
	return matched
}

func isValidGuestNumber(accompanyingGuests int) bool {
	return accompanyingGuests >= 0
}

func isValidTableNumber(tableNumber int, tableCount int) bool {
	return tableCount >= tableNumber && tableNumber > 0
}
