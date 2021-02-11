package system_validator

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"regexp"
)

type systemValidator struct {
	config conf.Configuration
}

func (sv systemValidator) IsValidGuestName(name string) bool {
	matched, err := regexp.MatchString(`[a-z,.'-]+`, name)
	if err != nil {
		return false
	}
	return matched
}

func (sv systemValidator) IsValidGuestNumber(accompanyingGuests int) bool {
	return accompanyingGuests >= 0
}

func (sv systemValidator) IsValidTableNumber(tableNumber int) bool {
	return sv.config.Tables.TableCount >= tableNumber && tableNumber > 0
}

func NewSystemValidator(config conf.Configuration) *systemValidator {
	return &systemValidator{
		config:  config,
	}
}
