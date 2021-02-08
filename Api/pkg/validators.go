package pkg

import (
	"regexp"
	"sync"
)

func (newGuest Guest) validate(tableCount int) bool {
	var wg *sync.WaitGroup

	defer wg.Done()
}


// general
func isValidGuestName(name string) bool{
	matched, err := regexp.MatchString(`[a-z ,.'-]+`, name)
	if err != nil {
		return false
	}
	return matched
}

func isValidGuestNumber(accompanyingGuests int) bool {
	return accompanyingGuests >= 0
}

func isValidTableNumber(tableNumber int, tableCount int) bool {
	return tableCount >= tableNumber
}