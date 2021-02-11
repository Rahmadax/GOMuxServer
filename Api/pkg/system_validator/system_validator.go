package system_validator

import (
	"errors"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"regexp"
)

type systemValidator struct {
	config conf.Configuration
}

func (sv systemValidator) ValidateNewGuest(newGuest models.Guest) error {
	err := sv.ValidateGuestName(newGuest.Name)
	if err != nil {
		return err
	}

	err = sv.ValidateTableNumber(newGuest.Table)
	if err != nil {
		return err
	}

	err = sv.ValidateAccompanyingGuests(newGuest.AccompanyingGuests)
	if err != nil {
		return err
	}

	return nil
}

func (sv systemValidator) ValidateArrivingGuest(guestName string, accompanyingGuests int) error {
	err := sv.ValidateGuestName(guestName)
	if err != nil {
		return err
	}

	err = sv.ValidateAccompanyingGuests(accompanyingGuests)
	if err != nil {
		return err
	}

	return nil
}

func (sv systemValidator) ValidateGuestName(name string) error {
	matched, err := regexp.MatchString(`[a-z,.'-]+`, name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid guest name")
	}

	return nil
}

func (sv systemValidator) ValidateAccompanyingGuests(accompanyingGuests int) error {
	if accompanyingGuests < 0 {
		return errors.New("guest can't have negative accompanying guests")
	}

	return nil
}

func (sv systemValidator) ValidateTableNumber(tableNumber int) error {
	if tableNumber > sv.config.Tables.TableCount {
		return errors.New(fmt.Sprintf("there are only %d tables", sv.config.Tables.TableCount))
	}

	if tableNumber < 0 {
		return errors.New("table number must be larger than 0")
	}

	return nil
}


func NewSystemValidator(config conf.Configuration) *systemValidator {
	return &systemValidator{
		config:  config,
	}
}
