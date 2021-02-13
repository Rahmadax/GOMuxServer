package empty_seats

import (
	"encoding/json"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
)

type EmptySeatsService interface {
	countPresentGuests() (int, error)
}

type emptySeatsHandler struct {
	config  conf.Configuration
	service EmptySeatsService
}

func NewEmptySeatsHandler(service EmptySeatsService, config conf.Configuration)  *emptySeatsHandler {
	return &emptySeatsHandler {
		config: config,
		service: service,
	}
}

func AddGuestListRoutes(config conf.Configuration, service EmptySeatsService, router *mux.Router){
	handler := NewEmptySeatsHandler(service,config)

	router.HandleFunc(config.Routes.CountEmptySeatsUri, handler.countEmptySeatsHandler()).Methods("GET")
}

func (esHandler *emptySeatsHandler) countEmptySeatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		presentGuestCount, err := esHandler.service.countPresentGuests()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, _ := json.Marshal(models.SeatsEmptyResponse{SeatsEmpty: esHandler.config.Tables.TotalCapacity - presentGuestCount})
		_, _ = w.Write(response)
	}
}