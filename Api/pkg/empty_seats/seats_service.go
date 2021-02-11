package empty_seats

import (
	"github.com/Rahmadax/GOMuxServer/Api/conf"
)

type GuestsRepository interface {
	CountPresentGuests() (int, error)
}

type emptySeatsService struct {
	config   conf.Configuration
	guestsRepo GuestsRepository
}

func (esService emptySeatsService) countPresentGuests() (int, error) {
	presentGuestCount, err := esService.guestsRepo.CountPresentGuests()
	if err != nil {
		return 0, err
	}

	return presentGuestCount, nil
}


func NewEmptySeatsService(config conf.Configuration, guestsRepo GuestsRepository) *emptySeatsService {
	return &emptySeatsService{
		config,
		guestsRepo,
	}
}
